package middleware

import (
	"net/http"
	"time"
)

// Timeout limits request handling time.
func Timeout(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, duration, "request timed out")
	}
}
