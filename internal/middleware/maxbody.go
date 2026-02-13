package middleware

import (
	"io"
	"net/http"
)

type limitedReadCloser struct {
	io.Reader
	io.Closer
}

func MaxBodyBytes(next http.Handler, maxBytes int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			next.ServeHTTP(w, r)
			return
		}
		r.Body = &limitedReadCloser{
			Reader: io.LimitReader(r.Body, maxBytes+1),
			Closer: r.Body,
		}
		next.ServeHTTP(w, r)
	})
}
