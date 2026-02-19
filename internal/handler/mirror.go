package handler

import (
	stderrors "errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"trykkeri-api/internal/errors"
)

const (
	mirrorFetchTimeout = 15 * time.Second
	maxURLBodyBytes    = 8192 // plenty for any URL
)

// Mirror fetches the HTML at the given URL (from request body) and renders it to PDF (same options as /print).
func (h *Handler) Mirror(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, maxURLBodyBytes+1))
	if err != nil {
		errors.WriteHTTP(r.Context(), w, errors.Internal("failed to read body: %v", err))
		return
	}
	if len(body) > maxURLBodyBytes {
		errors.WriteHTTP(r.Context(), w, errors.InvalidInput("url too long"))
		return
	}

	rawURL := strings.TrimSpace(string(body))
	if rawURL == "" {
		errors.WriteHTTP(r.Context(), w, errors.InvalidInput("request body must contain the URL"))
		return
	}

	targetURL, err := url.Parse(rawURL)
	if err != nil {
		errors.WriteHTTP(r.Context(), w, errors.InvalidInput("invalid url: %v", err))
		return
	}
	if err := validateMirrorRequestURL(targetURL); err != nil {
		errors.WriteHTTP(r.Context(), w, err)
		return
	}

	client := newMirrorHTTPClient()
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, targetURL.String(), nil)
	if err != nil {
		errors.WriteHTTP(r.Context(), w, errors.Internal("failed to create request: %v", err))
		return
	}
	req.Header.Set("User-Agent", "Trykkeri-API-Mirror/1.0")

	resp, err := client.Do(req)
	if err != nil {
		if stderrors.Is(err, errors.ErrInvalidInput) {
			errors.WriteHTTP(r.Context(), w, err)
			return
		}
		errors.WriteHTTP(r.Context(), w, errors.PdfGeneration("fetch failed: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errors.WriteHTTP(r.Context(), w, errors.PdfGeneration("fetch failed: %s", resp.Status))
		return
	}

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, h.cfg.MaxBodyBytes+1))
	if err != nil {
		errors.WriteHTTP(r.Context(), w, errors.Internal("failed to read response: %v", err))
		return
	}
	if int64(len(respBody)) > h.cfg.MaxBodyBytes {
		errors.WriteHTTP(r.Context(), w, errors.ErrPayloadTooLarge)
		return
	}

	html := string(respBody)
	if strings.TrimSpace(html) == "" {
		errors.WriteHTTP(r.Context(), w, errors.InvalidInput("target page returned empty content"))
		return
	}

	query := r.URL.Query()
	opts := queryToPdfOptions(query)
	baseURL := query.Get("base_url")
	var baseURLPtr *string
	if baseURL != "" {
		baseURLPtr = &baseURL
	} else {
		// Default base_url to the fetched page so relative links (CSS, images) resolve
		baseURLStr := targetURL.String()
		baseURLPtr = &baseURLStr
	}

	pdfBytes, err := h.pdfSvc.Render(r.Context(), html, baseURLPtr, opts)
	if err != nil {
		errors.WriteHTTP(r.Context(), w, err)
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
