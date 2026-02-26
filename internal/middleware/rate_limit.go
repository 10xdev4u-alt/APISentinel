package middleware

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/princetheprogrammer/apisentinel/internal/logger"
)

// RateLimiter simple in-memory rate limiter.
// In a real pro-tier app, you'd use Redis or a Token Bucket library.
type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]int
	limit   int
}

func NewRateLimiter(limit int) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]int),
		limit:   limit,
	}

	// Reset limits every minute (educational simplification)
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			rl.mu.Lock()
			rl.clients = make(map[string]int)
			rl.mu.Unlock()
		}
	}()

	return rl
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr // Fallback
		}

		rl.mu.Lock()
		rl.clients[ip]++
		count := rl.clients[ip]
		rl.mu.Unlock()

		if count > rl.limit {
			log.Printf("⚠️ Rate Limit Exceeded for IP: %s", ip)
			requestID := r.Header.Get("X-Request-ID")
			logger.LogEvent(requestID, ip, r.Method, r.URL.Path, "Rate Limit Exceeded", "Client exceeded allowed requests per minute")
			IncrementBlocked()
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
