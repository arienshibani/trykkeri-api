package pdf

import (
	"testing"
)

func TestDefaultPdfOptions(t *testing.T) {
	opts := DefaultPdfOptions()
	if opts.PageSize == nil || *opts.PageSize != "A4" {
		t.Errorf("PageSize = %v; want A4", opts.PageSize)
	}
	if opts.MarginTopMm == nil || *opts.MarginTopMm != 10 {
		t.Errorf("MarginTopMm = %v; want 10", opts.MarginTopMm)
	}
	if opts.DPI == nil || *opts.DPI != 300 {
		t.Errorf("DPI = %v; want 300", opts.DPI)
	}
	if opts.Portrait == nil || !*opts.Portrait {
		t.Errorf("Portrait = %v; want true", opts.Portrait)
	}
}
