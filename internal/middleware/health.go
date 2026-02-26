package middleware

import (
	"encoding/json"
	"net/http"
)

// HealthResponse represents the health status of the system.
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// HealthHandler returns a 200 OK status.
func HealthHandler(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthResponse{
			Status:  "UP",
			Version: version,
		})
	}
}
