package middleware

import "net/http"

// SecurityHeaders adds common security headers to all outgoing responses.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent browsers from MIME-sniffing the response
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Prevent Clickjacking attacks
		w.Header().Set("X-Frame-Options", "DENY")

		// Enable XSS filtering in browsers
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Content Security Policy (Basic)
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		// Pass to the next handler
		next.ServeHTTP(w, r)
	})
}
