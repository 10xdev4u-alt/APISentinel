package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/princetheprogrammer/apisentinel/internal/logger"
	"github.com/princetheprogrammer/apisentinel/internal/middleware"
	"github.com/princetheprogrammer/apisentinel/internal/proxy"
	"github.com/princetheprogrammer/apisentinel/internal/testserver"
)

func main() {
	// Initialize Audit Logger
	if err := logger.InitAuditLogger("audit.log"); err != nil {
		log.Fatalf("❌ Failed to initialize audit logger: %v", err)
	}
	defer logger.Close()

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

	// 2. Setup the Multi-Target Proxy
	mtProxy := proxy.NewMultiTargetProxy()

	// Default Route
	if err := mtProxy.AddRoute("/", backendURL); err != nil {
		log.Fatalf("❌ Failed to add default route: %v", err)
	}

	// Add a second route for demonstration
	apiV2URL := getEnv("API_V2_URL", "")
	if apiV2URL == "" {
		apiV2Port := "9001"
		apiV2URL = "http://localhost:" + apiV2Port
		log.Printf("📦 Starting second mock backend for /api/v2 on :%s", apiV2Port)
		go testserver.StartTestServer(apiV2Port)
	}
	if err := mtProxy.AddRoute("/api/v2", apiV2URL); err != nil {
		log.Fatalf("❌ Failed to add /api/v2 route: %v", err)
	}

	// Initialize Middlewares
	rl := middleware.NewRateLimiter(*rateLimit)
	inspector := middleware.NewSecurityInspector()
	adminKey := getEnv("ADMIN_KEY", "secret-sentinel-key")
	blocklist := middleware.NewIPBlocklist(adminKey)

	// 3. Start the API Sentinel Proxy Server
	log.Printf("🛡️ API Sentinel Proxy starting on :%s", proxyPort)
	log.Printf("🛡️ Forwarding traffic to: %s", backendURL)

	// Create a multiplexer to handle /stats separately
	mux := http.NewServeMux()
	mux.HandleFunc("/stats", middleware.StatsHandler)
	mux.HandleFunc("/block", blocklist.AdminHandler)
	mux.HandleFunc("/unblock", blocklist.AdminHandler)

	// Apply Middlewares to the Proxy
	proxyWithMiddleware := middleware.Chain(mtProxy,
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				middleware.IncrementTotal()
				next.ServeHTTP(w, r)
			})
		},
		blocklist.Middleware,
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
