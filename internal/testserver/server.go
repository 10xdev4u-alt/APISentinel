package testserver

import (
	"fmt"
	"net/http"
)

// StartTestServer starts a simple backend server for testing the proxy.
func StartTestServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "🚀 Hello from Backend Server on port %s!", port)
	})

	mux.HandleFunc("/leak", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Oops! Here is a credit card: 4111-1111-1111-1111")
	})

	fmt.Printf("📦 Test Backend Server starting on :%s...\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Printf("❌ Backend Server Error: %v\n", err)
	}
}
