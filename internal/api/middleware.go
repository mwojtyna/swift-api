package api

import (
	"log"
	"net/http"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

// Override
func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func LoggingMiddleware(next http.Handler, l *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapped, r)

		l.Printf("%d %s %s from %s", wrapped.statusCode, r.Method, r.URL.Path, r.RemoteAddr)
	})
}
