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

	// 3. Start the API Sentinel Proxy Server
	proxyPort := "8080"
	log.Printf("🛡️ API Sentinel Proxy starting on :%s", proxyPort)
	log.Printf("🛡️ Forwarding traffic to: %s", targetURL)

	// Apply Middlewares
	handler := middleware.Chain(proxyServer,
		rl.Middleware,
		middleware.SecurityHeaders,
	)

	if err := http.ListenAndServe(":"+proxyPort, handler); err != nil {
		log.Fatalf("❌ API Sentinel Proxy Error: %v", err)
	}
}
