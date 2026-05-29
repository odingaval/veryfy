package middleware

import (
	"log/slog"
	"net/http"

	"github.com/odingaval/veryfy/api/internal/httpjson"
)

// Recover converts unexpected panics into JSON 500 responses.
func Recover(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if value := recover(); value != nil {
					logger.Error("panic recovered", "panic", value, "method", r.Method, "path", r.URL.Path)
					httpjson.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
