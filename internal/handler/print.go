package handler

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"trykkeri-api/internal/errors"
	"trykkeri-api/internal/middleware"
	"trykkeri-api/internal/pdf"
)

func (h *Handler) Print(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	body, err := io.ReadAll(io.LimitReader(r.Body, h.cfg.MaxBodyBytes+1))
	if err != nil {
		errors.WriteHTTP(w, errors.Internal("failed to read body: %v", err))
		return
	}
	if int64(len(body)) > h.cfg.MaxBodyBytes {
		errors.WriteHTTP(w, errors.ErrPayloadTooLarge)
		return
	}

	html := string(body)
	if strings.TrimSpace(html) == "" {
		errors.WriteHTTP(w, errors.InvalidInput("HTML content cannot be empty"))
		return
	}

	if max := h.cfg.PayloadLogMaxBytes; max > 0 {
		preview := html
		if len(html) > max {
			preview = html[:max]
		}
		middleware.AddRequestLogAttrs(r.Context(), "payload_preview", preview, "payload_size", len(html))
	}

	opts := queryToPdfOptions(query)
	baseURL := query.Get("base_url")
	var baseURLPtr *string
	if baseURL != "" {
		baseURLPtr = &baseURL
	}

	pdfBytes, err := h.pdfSvc.Render(r.Context(), html, baseURLPtr, opts)
	if err != nil {
		errors.WriteHTTP(w, err)
		return
	}

	filename := query.Get("filename")
	if filename == "" {
		filename = "document.pdf"
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `inline; filename="`+filename+`"`)
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pdfBytes)
}

func queryToPdfOptions(q url.Values) *pdf.PdfOptions {
	getStr := func(key string) string {
		return q.Get(key)
	}
	getUint32 := func(key string) *uint32 {
		s := getStr(key)
		if s == "" {
			return nil
		}
		n, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return nil
		}
		u := uint32(n)
		return &u
	}
	getBool := func(key string) *bool {
		s := getStr(key)
		if s == "" {
			return nil
		}
		v, err := strconv.ParseBool(s)
		if err != nil {
			return nil
		}
		return &v
	}

	pageSize := getStr("page_size")
	portraitSet := getStr("portrait") != ""
	if pageSize == "" && !portraitSet &&
		getStr("margin_top_mm") == "" && getStr("margin_right_mm") == "" &&
		getStr("margin_bottom_mm") == "" && getStr("margin_left_mm") == "" &&
		getStr("dpi") == "" && getStr("print_background") == "" && getStr("grayscale") == "" {
		return nil
	}

	opts := pdf.DefaultPdfOptions()
	if pageSize != "" {
		opts.PageSize = &pageSize
	}
	if v := getBool("portrait"); v != nil {
		opts.Portrait = v
	}
	if v := getUint32("margin_top_mm"); v != nil {
		opts.MarginTopMm = v
	}
	if v := getUint32("margin_right_mm"); v != nil {
		opts.MarginRightMm = v
	}
	if v := getUint32("margin_bottom_mm"); v != nil {
		opts.MarginBottomMm = v
	}
	if v := getUint32("margin_left_mm"); v != nil {
		opts.MarginLeftMm = v
	}
	if v := getUint32("dpi"); v != nil {
		opts.DPI = v
	}
	if v := getBool("print_background"); v != nil {
		opts.PrintBackground = v
	}
	if v := getBool("grayscale"); v != nil {
		opts.Grayscale = v
	}
	return &opts
}
