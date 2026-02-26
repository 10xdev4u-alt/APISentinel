package middleware

import (
	"bytes"
	"log"
	"net/http"
	"regexp"

	"github.com/princetheprogrammer/apisentinel/internal/logger"
)

// DLPMiddleware scans outgoing responses for sensitive data leaks.
type DLPMiddleware struct {
	patterns []Pattern
}

func NewDLPMiddleware() *DLPMiddleware {
	return &DLPMiddleware{
		patterns: []Pattern{
			{
				Name:    "Credit Card Leak",
				Regexp:  regexp.MustCompile(`\b(?:\d[ -]*?){13,16}\b`),
				Message: "Sensitive financial data detected in response.",
			},
			{
				Name:    "SSN Leak",
				Regexp:  regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`),
				Message: "Personally Identifiable Information (SSN) detected in response.",
			},
		},
	}
}

// bodyLogWriter wraps http.ResponseWriter to capture the response body and status code.
type bodyLogWriter struct {
	http.ResponseWriter
	body           *bytes.Buffer
	statusCode     int
	headerWritten  bool
}

func (w *bodyLogWriter) WriteHeader(code int) {
	if !w.headerWritten {
		w.statusCode = code
		// We DON'T call w.ResponseWriter.WriteHeader(code) yet!
		w.headerWritten = true
	}
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (dlp *DLPMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blw := &bodyLogWriter{
			body:           bytes.NewBuffer(nil),
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default
		}

		next.ServeHTTP(blw, r)

		responseBody := blw.body.Bytes()
		for _, p := range dlp.patterns {
			if p.Regexp.Match(responseBody) {
				log.Printf("🛑 DLP ALERT [%s]: Blocked outgoing response to %s", p.Name, r.RemoteAddr)
				requestID := r.Header.Get("X-Request-ID")
				logger.LogEvent(requestID, r.RemoteAddr, r.Method, r.URL.Path, "DLP Violation: "+p.Name, "Backend attempted to leak sensitive data")

				w.Header().Del("Content-Length")
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.Header().Set("X-Sentinel-DLP", "Blocked")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Security Error: Data Loss Prevention policy triggered. Response blocked."))
				return
			}
		}

		// If we reach here, it's safe. Write the original status and body.
		if blw.headerWritten {
			w.WriteHeader(blw.statusCode)
		}
		w.Write(responseBody)
	})
}
