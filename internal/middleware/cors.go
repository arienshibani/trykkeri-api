package middleware

import (
	"net/http"
)

func CORS(next http.Handler, allowedOrigins []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			allow := len(allowedOrigins) == 0
			if !allow {
				for _, o := range allowedOrigins {
					if o == origin || o == "*" {
						allow = true
						break
					}
				}
			}
			if allow {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
