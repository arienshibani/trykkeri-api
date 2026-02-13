package pdf

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"trykkeri-api/internal/config"
	"trykkeri-api/internal/errors"
)

type PdfOptions struct {
	PageSize        *string
	MarginTopMm     *uint32
	MarginRightMm   *uint32
	MarginBottomMm  *uint32
	MarginLeftMm    *uint32
	DPI             *uint32
	PrintBackground *bool
	Grayscale       *bool
	Portrait        *bool // true = portrait, false = landscape
}

func DefaultPdfOptions() PdfOptions {
	a4 := "A4"
	var ten uint32 = 10
	var dpiDefault uint32 = 300
	trueVal := true
	falseVal := false
	return PdfOptions{
		PageSize:        &a4,
		MarginTopMm:     &ten,
		MarginRightMm:   &ten,
		MarginBottomMm:  &ten,
		MarginLeftMm:    &ten,
		DPI:             &dpiDefault,
		PrintBackground: &trueVal,
		Grayscale:       &falseVal,
		Portrait:        &trueVal,
	}
}

type Service struct {
	cfg *config.Config
}

func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

func (s *Service) Render(ctx context.Context, html string, baseURL *string, opts *PdfOptions) ([]byte, error) {
	if opts == nil {
		def := DefaultPdfOptions()
		opts = &def
	}

	dir, err := os.MkdirTemp("", "trykkeri-api-*")
	if err != nil {
		return nil, errors.Internal("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	inputPath := filepath.Join(dir, "input.html")
	if err := os.WriteFile(inputPath, []byte(html), 0644); err != nil {
		return nil, errors.Internal("failed to write HTML: %v", err)
	}

	outputPath := filepath.Join(dir, "output.pdf")

	args := []string{"--quiet", "--encoding", "utf-8"}

	if opts.PageSize != nil {
		args = append(args, "--page-size", *opts.PageSize)
	}
	if opts.DPI != nil {
		args = append(args, "--dpi", fmt.Sprintf("%d", *opts.DPI))
	}
	if opts.Portrait == nil || *opts.Portrait {
		args = append(args, "--orientation", "Portrait")
	} else {
		args = append(args, "--orientation", "Landscape")
	}
	if opts.MarginTopMm != nil {
		args = append(args, "--margin-top", fmt.Sprintf("%dmm", *opts.MarginTopMm))
	}
	if opts.MarginRightMm != nil {
		args = append(args, "--margin-right", fmt.Sprintf("%dmm", *opts.MarginRightMm))
	}
	if opts.MarginBottomMm != nil {
		args = append(args, "--margin-bottom", fmt.Sprintf("%dmm", *opts.MarginBottomMm))
	}
	if opts.MarginLeftMm != nil {
		args = append(args, "--margin-left", fmt.Sprintf("%dmm", *opts.MarginLeftMm))
	}
	if opts.PrintBackground == nil || *opts.PrintBackground {
		args = append(args, "--print-media-type")
	}
	if opts.Grayscale != nil && *opts.Grayscale {
		args = append(args, "--grayscale")
	}

	if !s.cfg.AllowNet {
		args = append(args, "--disable-external-links")
	}
	for _, p := range s.cfg.AllowlistPaths {
		args = append(args, "--allow", p)
	}

	args = append(args, inputPath, outputPath)

	timeoutDur := time.Duration(s.cfg.RenderTimeoutMs) * time.Millisecond
	runCtx, cancel := context.WithTimeout(ctx, timeoutDur)
	defer cancel()

	cmd := exec.CommandContext(runCtx, s.cfg.WkhtmltopdfPath, args...)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		if runCtx.Err() == context.DeadlineExceeded {
			return nil, errors.ErrTimeout
		}
		return nil, errors.PdfGeneration("wkhtmltopdf failed: %s", string(out))
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, errors.Internal("failed to read PDF output: %v", err)
	}
	if len(data) == 0 {
		return nil, errors.PdfGeneration("generated PDF is empty")
	}
	return data, nil
}
