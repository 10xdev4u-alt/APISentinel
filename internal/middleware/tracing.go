package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// Tracing adds a unique X-Request-ID to each request.
func Tracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// Inject into Request headers (for the backend)
		r.Header.Set("X-Request-ID", requestID)

		// Inject into Response headers (for the client)
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r)
	})
}

func generateRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}
