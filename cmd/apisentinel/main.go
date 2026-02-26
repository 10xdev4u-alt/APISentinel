package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/princetheprogrammer/apisentinel/internal/config"
	"github.com/princetheprogrammer/apisentinel/internal/logger"
	"github.com/princetheprogrammer/apisentinel/internal/middleware"
	"github.com/princetheprogrammer/apisentinel/internal/proxy"
	"github.com/princetheprogrammer/apisentinel/internal/testserver"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	// 1. Load Configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Printf("⚠️ Config file not found or invalid (%s). Using defaults.", *configPath)
		cfg = &config.Config{
			Server: config.ServerConfig{Port: 8080, AdminKey: "secret-sentinel-key", RateLimit: 10, AuditLog: "audit.log"},
			Security: config.SecurityConfig{EnableXSS: true, EnableSQLi: true, EnableDLP: true},
		}
	}

	// 2. Initialize Audit Logger
	if err := logger.InitAuditLogger(cfg.Server.AuditLog); err != nil {
		log.Fatalf("❌ Failed to initialize audit logger: %v", err)
	}
	defer logger.Close()

	// 3. Setup the Multi-Target Proxy
	mtProxy := proxy.NewMultiTargetProxy()
	if len(cfg.Routes) == 0 {
		log.Println("📦 No routes in config. Starting mock backends.")
		// Default Mock Setup
		p1, p2 := "9000", "9001"
		go testserver.StartTestServer(p1)
		go testserver.StartTestServer(p2)
		mtProxy.AddRoute("/", []string{"http://localhost:" + p1})
		mtProxy.AddRoute("/api/v2", []string{"http://localhost:" + p2})
	} else {
		for _, r := range cfg.Routes {
			targets := r.Targets
			if len(targets) == 0 && r.Target != "" {
				targets = []string{r.Target}
			}
			
			log.Printf("🛣️  Adding route: %s -> %v", r.Path, targets)
			if err := mtProxy.AddRoute(r.Path, targets); err != nil {
				log.Fatalf("❌ Failed to add route %s: %v", r.Path, err)
			}
		}
	}

	// 4. Initialize Middlewares
	rl := middleware.NewRateLimiter(cfg.Server.RateLimit)
	inspector := middleware.NewSecurityInspector(cfg.Security.EnableXSS, cfg.Security.EnableSQLi)
	blocklist := middleware.NewIPBlocklist(cfg.Server.AdminKey)
	
	var dlp *middleware.DLPMiddleware
	if cfg.Security.EnableDLP {
		dlp = middleware.NewDLPMiddleware()
	}

	// 5. Start the API Sentinel Proxy Server
	proxyPort := fmt.Sprintf("%d", cfg.Server.Port)
	log.Printf("🛡️ API Sentinel Proxy starting on :%s", proxyPort)

	// Create a multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/health", middleware.HealthHandler("1.0.0"))
	mux.HandleFunc("/stats", middleware.StatsHandler)
	mux.HandleFunc("/block", blocklist.AdminHandler)
	mux.HandleFunc("/unblock", blocklist.AdminHandler)
	mux.HandleFunc("/dashboard", middleware.DashboardHandler(blocklist, cfg.Server.AuditLog))

	// Build the middleware chain
	mws := []middleware.Middleware{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				middleware.IncrementTotal()
				next.ServeHTTP(w, r)
			})
		},
		blocklist.Middleware,
	}

	if dlp != nil {
		mws = append(mws, dlp.Middleware)
	}

	mws = append(mws,
		inspector.Middleware,
		rl.Middleware,
		middleware.SecurityHeaders,
	)

	// Apply Middlewares to the Proxy
	proxyWithMiddleware := middleware.Chain(mtProxy, mws...)

	// Route everything else to the proxy
	mux.Handle("/", proxyWithMiddleware)

	// --- 6. Start Server with Graceful Shutdown ---
	server := &http.Server{
		Addr:    ":" + proxyPort,
		Handler: mux,
	}

	// Channel to listen for interrupt signals (Ctrl+C, SIGTERM)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run server in a goroutine
	go func() {
		log.Printf("🛡️ API Sentinel Proxy starting on :%s", proxyPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ API Sentinel Proxy Error: %v", err)
		}
	}()

	// Wait for signal
	<-stop
	log.Println("\n🛑 Shutdown signal received. Cleaning up...")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server Shutdown Failed: %v", err)
	}

	log.Println("✅ API Sentinel Shutdown Gracefully. See you next time, Prince!")
}
