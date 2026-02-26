package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/princetheprogrammer/apisentinel/internal/logger"
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

func NewSecurityInspector(enableXSS, enableSQLi bool) *SecurityInspector {
	si := &SecurityInspector{}

	if enableXSS {
		si.patterns = append(si.patterns, Pattern{
			Name:    "XSS Detection",
			Regexp:  regexp.MustCompile(`(?i)<script.*?>|javascript:|onload=`),
			Message: "Malicious <script> or javascript: detected.",
		})
	}

	if enableSQLi {
		si.patterns = append(si.patterns, Pattern{
			Name:    "SQL Injection Detection",
			Regexp:  regexp.MustCompile(`(?i)(union.*select|insert.*into|drop.*table|truncate.*table|' or 1=1|--|#)`),
			Message: "SQL injection attempt detected.",
		})
	}

	return si
}

func (si *SecurityInspector) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Inspect Query Parameters
		if matched, pattern := si.inspectQuery(r); matched {
			si.block(w, r, "Query", pattern)
			return
		}

		// 2. Inspect Request Body (if any)
		if si.inspectBody(w, r) {
			// inspectBody handles the blocking and error responses internally if it returns true
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (si *SecurityInspector) inspectQuery(r *http.Request) (bool, string) {
	return si.isMalicious(r.URL.RawQuery)
}

func (si *SecurityInspector) inspectBody(w http.ResponseWriter, r *http.Request) bool {
	if r.Body == nil || r.ContentLength == 0 {
		return false
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("❌ Failed to read request body: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return true // Stop processing on error
	}

	// Restore the body for the next handler
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	contentType := r.Header.Get("Content-Type")

	// 1. Check for JSON
	if strings.Contains(contentType, "application/json") {
		var jsonData interface{}
		if err := json.Unmarshal(bodyBytes, &jsonData); err == nil {
			if matched, pattern := si.inspectValue(jsonData); matched {
				si.block(w, r, "JSON Body", pattern)
				return true
			}
			return false
		}
	}

	// 2. Check for Form Data
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		values, err := url.ParseQuery(string(bodyBytes))
		if err == nil {
			for _, v := range values {
				for _, val := range v {
					if matched, pattern := si.isMalicious(val); matched {
						si.block(w, r, "Form Body", pattern)
						return true
					}
				}
			}
			return false
		}
	}

	// 3. Default String Match for other body types
	if matched, pattern := si.isMalicious(string(bodyBytes)); matched {
		si.block(w, r, "Body", pattern)
		return true
	}

	return false
}

func (si *SecurityInspector) inspectValue(v interface{}) (bool, string) {
	switch val := v.(type) {
	case string:
		return si.isMalicious(val)
	case []interface{}:
		for _, item := range val {
			if matched, pattern := si.inspectValue(item); matched {
				return true, pattern
			}
		}
	case map[string]interface{}:
		for k, v := range val {
			if matched, pattern := si.isMalicious(k); matched {
				return true, pattern
			}
			if matched, pattern := si.inspectValue(v); matched {
				return true, pattern
			}
		}
	}
	return false, ""
}

func (si *SecurityInspector) isMalicious(data string) (bool, string) {
	if data == "" {
		return false, ""
	}

	// Decode URL-encoded strings for inspection
	decoded, err := url.QueryUnescape(data)
	if err != nil {
		decoded = data // Fallback to original
	}

	for _, p := range si.patterns {
		if p.Regexp.MatchString(decoded) {
			log.Printf("🛑 BLOCKED [%s]: Found in input", p.Name)
			return true, p.Name
		}
	}
	return false, ""
}

func (si *SecurityInspector) block(w http.ResponseWriter, r *http.Request, source, pattern string) {
	log.Printf("🛡️ API Sentinel: Blocking request from %s due to malicious content", source)
	
	// Get IP
	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = forwarded
	}

	// Log to Audit File
	requestID := r.Header.Get("X-Request-ID")
	logger.LogEvent(requestID, ip, r.Method, r.URL.Path, pattern, "Blocked in: "+source)

	IncrementBlocked()
	http.Error(w, "Forbidden: Malicious activity detected", http.StatusForbidden)
}
