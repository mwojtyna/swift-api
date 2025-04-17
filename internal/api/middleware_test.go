package api

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.Handler
		path           string
		addr           string
		expectedStatus int
		expectedLog    string
	}{
		{
			name: "200 ok",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			path:           "/test",
			addr:           "1.2.3.4",
			expectedStatus: http.StatusOK,
			expectedLog:    "200 GET /test from 1.2.3.4",
		},
		{
			name: "404 not found",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}),
			path:           "/notfound",
			addr:           "5.6.7.8",
			expectedStatus: http.StatusNotFound,
			expectedLog:    "404 GET /notfound from 5.6.7.8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logBuf bytes.Buffer
			logger := log.New(&logBuf, "", 0)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			req.RemoteAddr = tt.addr

			w := httptest.NewRecorder()

			middleware := LoggingMiddleware(tt.handler, logger)
			middleware.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, logBuf.String(), tt.expectedLog)
		})
	}
}
