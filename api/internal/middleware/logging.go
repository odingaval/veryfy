package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// RequestLogger records a compact log line for each request.
func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startedAt := time.Now()

			next.ServeHTTP(w, r)

			logger.Info("request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"durationMs", time.Since(startedAt).Milliseconds(),
			)
		})
	}
}
