package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status        string `json:"status"`
	Version       string `json:"version"`
	UptimeSeconds int64  `json:"uptime_seconds"`
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.startTime).Seconds()
	resp := HealthResponse{
		Status:        "ok",
		Version:       h.version,
		UptimeSeconds: int64(uptime),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if r.Method != http.MethodHead {
		_ = json.NewEncoder(w).Encode(resp)
	}
}

// Favicon responds to browser favicon requests with 204 No Content so they don't 404.
func (h *Handler) Favicon(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
