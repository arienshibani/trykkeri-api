package errors

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteHTTP(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
		wantBody string
	}{
		{"invalid input", InvalidInput("bad"), http.StatusBadRequest, "invalid_input"},
		{"timeout", ErrTimeout, http.StatusRequestTimeout, "timeout"},
		{"payload too large", ErrPayloadTooLarge, http.StatusRequestEntityTooLarge, "payload_too_large"},
		{"pdf generation", PdfGeneration("wk failed"), http.StatusInternalServerError, "pdf_generation_failed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteHTTP(w, tt.err)
			if w.Code != tt.wantCode {
				t.Errorf("Code = %d; want %d", w.Code, tt.wantCode)
			}
			body := w.Body.String()
			if body == "" || !strings.Contains(body, tt.wantBody) {
				t.Errorf("Body = %q; want to contain %q", body, tt.wantBody)
			}
		})
	}
}

