package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.statusCode = code
	sr.ResponseWriter.WriteHeader(code)
}

// requestLogging is a middleware that logs HTTP request info
func requestLogging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			statusWriter := &statusRecorder{w, http.StatusOK}
			next.ServeHTTP(statusWriter, r)

			logger.InfoContext(r.Context(),
				fmt.Sprintf("%d %s - %s", statusWriter.statusCode, http.StatusText(statusWriter.statusCode), r.URL.Path),
				"status", statusWriter.statusCode,
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent())
		})
	}
}
