package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/princetheprogrammer/apisentinel/internal/middleware"
	"github.com/princetheprogrammer/apisentinel/internal/proxy"
	"github.com/princetheprogrammer/apisentinel/internal/testserver"
)

func main() {
	// Flags and Env Vars
	proxyPort := getEnv("PROXY_PORT", "8080")
	backendURL := getEnv("BACKEND_URL", "")
	rateLimit := flag.Int("limit", 10, "Requests per minute")
	flag.Parse()

	// If no backend is provided, start the mock one for convenience
	if backendURL == "" {
		backendPort := "9000"
		backendURL = "http://localhost:" + backendPort
		log.Printf("📦 No BACKEND_URL provided. Starting mock backend on :%s", backendPort)
		go testserver.StartTestServer(backendPort)
	}

	// 2. Setup the Proxy
	proxyServer, err := proxy.NewProxy(backendURL)
	if err != nil {
		log.Fatalf("❌ Failed to create proxy: %v", err)
	}

	// Initialize Middlewares
	rl := middleware.NewRateLimiter(*rateLimit)
	inspector := middleware.NewSecurityInspector()

	// 3. Start the API Sentinel Proxy Server
	log.Printf("🛡️ API Sentinel Proxy starting on :%s", proxyPort)
	log.Printf("🛡️ Forwarding traffic to: %s", backendURL)

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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
