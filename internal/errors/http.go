package errors

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	stderrors "errors"

	"trykkeri-api/internal/middleware"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func WriteHTTP(ctx context.Context, w http.ResponseWriter, err error) {
	var status int
	var code, message string

	switch {
	case stderrors.Is(err, ErrInvalidInput):
		status = http.StatusBadRequest
		code = "invalid_input"
		message = err.Error()
	case stderrors.Is(err, ErrPdfGeneration):
		status = http.StatusInternalServerError
		code = "pdf_generation_failed"
		message = "PDF generation failed"
		slog.Error("PDF generation error", "err", err)
	case stderrors.Is(err, ErrTimeout):
		status = http.StatusRequestTimeout
		code = "timeout"
		message = "Request timeout"
	case stderrors.Is(err, ErrPayloadTooLarge):
		status = http.StatusRequestEntityTooLarge
		code = "payload_too_large"
		message = "Request body too large"
	default:
		status = http.StatusInternalServerError
		code = "internal_error"
		message = "Internal server error"
		slog.Error("Internal error", "err", err)
	}

	middleware.AddRequestLogAttrs(ctx, "error", err.Error())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: code, Message: message})
}
