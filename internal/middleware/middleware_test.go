package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddRequestLogAttrs_noOpWithoutMiddleware(t *testing.T) {
	// Should not panic when ctx has no log attrs key (e.g. outside RequestLog).
	AddRequestLogAttrs(context.Background(), "key", "value")
}

func TestRequestLog(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := RequestLog(next, "1.0.0")
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d; want 200", rec.Code)
	}
}
