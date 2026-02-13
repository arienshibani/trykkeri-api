package handler

import (
	"github.com/go-chi/chi/v5"
)

func Routes(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", h.Health)
	r.Head("/health", h.Health)
	r.Get("/favicon.ico", h.Favicon)
	r.Post("/print", h.Print)
	r.Post("/mirror", h.Mirror)
	r.Get("/openapi.json", h.OpenAPI)
	r.Get("/*", h.DocsUI)
	return r
}
