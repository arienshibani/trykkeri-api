package handler

import (
	_ "embed"
	"net/http"
)

//go:embed openapi.json
var openAPISpec []byte

//go:embed scalar.html
var docsHTML []byte

func (h *Handler) OpenAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(openAPISpec)
}

func (h *Handler) DocsUI(w http.ResponseWriter, r *http.Request) {
	// Only serve docs at exact "/" to avoid catching all routes
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(docsHTML)
}
