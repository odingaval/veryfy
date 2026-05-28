package handlers

import (
	"net/http"
	"time"

	"github.com/aochuka/veryfy/api/internal/httpjson"
)

// HealthHandler serves the health check endpoint.
type HealthHandler struct {
	startedAt time.Time
}

// NewHealthHandler creates a health handler.
func NewHealthHandler(startedAt time.Time) *HealthHandler {
	return &HealthHandler{
		startedAt: startedAt,
	}
}

// ServeHTTP returns a small readiness payload for load balancers and local checks.
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpjson.WriteSuccess(w, http.StatusOK, map[string]any{
		"status":    "ok",
		"service":   "veryfy-api",
		"uptimeSec": int(time.Since(h.startedAt).Seconds()),
	})
}
