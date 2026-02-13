package middleware

import (
	"net/http"
	"time"

	"trykkeri-api/internal/config"
)

func Chain(next http.Handler, cfg *config.Config, version string) http.Handler {
	next = Timeout(next, time.Duration(cfg.RenderTimeoutMs+5000)*time.Millisecond)
	next = MaxBodyBytes(next, cfg.MaxBodyBytes)
	next = Gzip(next)
	next = CORS(next, cfg.CORSOrigins)
	next = RequestLog(next, version)
	return next
}
