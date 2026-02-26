package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

// Pattern defines a malicious signature to look for.
type Pattern struct {
	Name    string
	Regexp  *regexp.Regexp
	Message string
}

// SecurityInspector scans requests for malicious patterns.
type SecurityInspector struct {
	patterns []Pattern
}

func NewSecurityInspector() *SecurityInspector {
	return &SecurityInspector{
		patterns: []Pattern{
			{
				Name:    "XSS Detection",
				Regexp:  regexp.MustCompile(`(?i)<script.*?>|javascript:|onload=`),
				Message: "Malicious <script> or javascript: detected.",
			},
			{
				Name:    "SQL Injection Detection",
				Regexp:  regexp.MustCompile(`(?i)(union.*select|insert.*into|drop.*table|truncate.*table|' or 1=1|--|#)`),
				Message: "SQL injection attempt detected.",
			},
		},
	}
}

func (si *SecurityInspector) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Inspect Query Parameters
		if si.isMalicious(r.URL.RawQuery) {
			si.block(w, "Query")
			return
		}

		// 2. Inspect Request Body (if any)
		if r.Body != nil && r.ContentLength > 0 {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("❌ Failed to read request body: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Restore the body for the next handler
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			if si.isMalicious(string(body)) {
				si.block(w, "Body")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (si *SecurityInspector) isMalicious(data string) bool {
	if data == "" {
		return false
	}

	// Decode URL-encoded strings for inspection
	decoded, err := url.QueryUnescape(data)
	if err != nil {
		decoded = data // Fallback to original
	}

	for _, p := range si.patterns {
		if p.Regexp.MatchString(decoded) {
			log.Printf("🛑 BLOCKED [%s]: Found in input", p.Name)
			return true
		}
	}
	return false
}

func (si *SecurityInspector) block(w http.ResponseWriter, source string) {
	log.Printf("🛡️ API Sentinel: Blocking request from %s due to malicious content", source)
	http.Error(w, "Forbidden: Malicious activity detected", http.StatusForbidden)
}
