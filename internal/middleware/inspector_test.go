package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSecurityInspector(t *testing.T) {
	inspector := NewSecurityInspector()

	tests := []struct {
		name           string
		url            string
		body           string
		expectedStatus int
	}{
		{
			name:           "Safe Request",
			url:            "/safe",
			body:           "hello world",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "XSS in Query",
			url:            "/?q=%3Cscript%3Ealert(1)%3C/script%3E",
			body:           "",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "SQLi in Query",
			url:            "/?user=admin%27%20OR%201%3D1%20--",
			body:           "",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "XSS in Body",
			url:            "/post",
			body:           "<script>alert(1)</script>",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "SQLi in Body",
			url:            "/login",
			body:           "username=admin' OR 1=1 --",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Nested JSON XSS",
			url:            "/api/data",
			body:           `{"user": {"name": "Prince", "bio": "<script>alert(1)</script>"}}`,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Form Data SQLi",
			url:            "/submit",
			body:           "username=admin&comment=drop table users--",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.url, strings.NewReader(tt.body))
			rr := httptest.NewRecorder()

			// Dummy handler that just returns 200 OK
			handler := inspector.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("%s: expected status %v, got %v", tt.name, tt.expectedStatus, rr.Code)
			}
		})
	}
}
