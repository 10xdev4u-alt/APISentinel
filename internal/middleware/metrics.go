package middleware

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
)

// Metrics keeps track of proxy statistics.
type Metrics struct {
	TotalRequests   uint64 `json:"total_requests"`
	BlockedRequests uint64 `json:"blocked_requests"`
}

var GlobalMetrics = &Metrics{}

func IncrementTotal() {
	atomic.AddUint64(&GlobalMetrics.TotalRequests, 1)
}

func IncrementBlocked() {
	atomic.AddUint64(&GlobalMetrics.BlockedRequests, 1)
}

// StatsHandler returns the current metrics as JSON.
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GlobalMetrics)
}
