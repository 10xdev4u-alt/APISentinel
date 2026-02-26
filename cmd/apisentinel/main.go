package main

import (
	"log"
	"net/http"

	"github.com/princetheprogrammer/apisentinel/internal/middleware"
	"github.com/princetheprogrammer/apisentinel/internal/proxy"
	"github.com/princetheprogrammer/apisentinel/internal/testserver"
)

func main() {
	// 1. Start a mock backend server in a goroutine
	backendPort := "9000"
	go testserver.StartTestServer(backendPort)

	// 2. Setup the Proxy
	targetURL := "http://localhost:" + backendPort
	proxyServer, err := proxy.NewProxy(targetURL)
	if err != nil {
		log.Fatalf("❌ Failed to create proxy: %v", err)
	}

	// Initialize Middlewares
	rl := middleware.NewRateLimiter(10) // 10 requests per minute
	inspector := middleware.NewSecurityInspector()

	// 3. Start the API Sentinel Proxy Server
	proxyPort := "8080"
	log.Printf("🛡️ API Sentinel Proxy starting on :%s", proxyPort)
	log.Printf("🛡️ Forwarding traffic to: %s", targetURL)

	// Create a multiplexer to handle /stats separately
	mux := http.NewServeMux()
	mux.HandleFunc("/stats", middleware.StatsHandler)

	// Apply Middlewares to the Proxy
	proxyWithMiddleware := middleware.Chain(proxyServer,
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				middleware.IncrementTotal()
				next.ServeHTTP(w, r)
			})
		},
		inspector.Middleware,
		rl.Middleware,
		middleware.SecurityHeaders,
	)

	// Route everything else to the proxy
	mux.Handle("/", proxyWithMiddleware)

	if err := http.ListenAndServe(":"+proxyPort, mux); err != nil {
		log.Fatalf("❌ API Sentinel Proxy Error: %v", err)
	}
}
