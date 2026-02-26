package middleware

import (
	"log"
	"net"
	"net/http"
	"sync"
)

// IPBlocklist manages a set of blocked IP addresses.
type IPBlocklist struct {
	mu       sync.RWMutex
	blocked  map[string]bool
	adminKey string
}

func NewIPBlocklist(adminKey string) *IPBlocklist {
	return &IPBlocklist{
		blocked:  make(map[string]bool),
		adminKey: adminKey,
	}
}

func (bl *IPBlocklist) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		bl.mu.RLock()
		isBlocked := bl.blocked[ip]
		bl.mu.RUnlock()

		if isBlocked {
			log.Printf("🚫 Blocked request from blacklisted IP: %s", ip)
			IncrementBlocked()
			http.Error(w, "Forbidden: Your IP is blacklisted", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AdminHandler handles /block and /unblock requests.
// In a real app, this would be a more robust API.
func (bl *IPBlocklist) AdminHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" || key != bl.adminKey {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ip := r.URL.Query().Get("ip")
	if ip == "" {
		http.Error(w, "Missing IP", http.StatusBadRequest)
		return
	}

	bl.mu.Lock()
	defer bl.mu.Unlock()

	if r.URL.Path == "/block" {
		bl.blocked[ip] = true
		log.Printf("🛡️ IP %s added to blocklist", ip)
		w.Write([]byte("IP Blocked"))
	} else if r.URL.Path == "/unblock" {
		delete(bl.blocked, ip)
		log.Printf("🔓 IP %s removed from blocklist", ip)
		w.Write([]byte("IP Unblocked"))
	}
}
