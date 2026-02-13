package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type logAttrsKey struct{}

// AddRequestLogAttrs appends key-value pairs to the current request's log line.
// Only has effect when called from a handler during RequestLog middleware.
func AddRequestLogAttrs(ctx context.Context, attrs ...any) {
	if v := ctx.Value(logAttrsKey{}); v != nil {
		if slice, ok := v.(*[]any); ok {
			*slice = append(*slice, attrs...)
		}
	}
}

func RequestLog(next http.Handler, version string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var extra []any
		r = r.WithContext(context.WithValue(r.Context(), logAttrsKey{}, &extra))
		ww := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(ww, r)
		attrs := []any{
			"method", r.Method,
			"uri", r.URL.Path,
			"status", ww.status,
			"duration_ms", time.Since(start).Milliseconds(),
		}
		attrs = append(attrs, extra...)
		if ww.status >= 400 {
			slog.Error("request", attrs...)
		} else {
			slog.Info("request", attrs...)
		}
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
