package handler

import (
	"time"

	"trykkeri-api/internal/config"
	"trykkeri-api/internal/pdf"
)

type Handler struct {
	cfg       *config.Config
	pdfSvc    *pdf.Service
	version   string
	startTime time.Time
}

func New(cfg *config.Config, pdfSvc *pdf.Service, version string, startTime time.Time) *Handler {
	return &Handler{
		cfg:       cfg,
		pdfSvc:    pdfSvc,
		version:   version,
		startTime: startTime,
	}
}
