package proxy

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// Target represents a single backend server.
type Target struct {
	URL   *url.URL
	Proxy *httputil.ReverseProxy
	Alive bool
	mu    sync.RWMutex
}

func (t *Target) SetAlive(alive bool) {
	t.mu.Lock()
	t.Alive = alive
	t.mu.Unlock()
}

func (t *Target) IsAlive() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.Alive
}

// LoadBalancer manages multiple targets for a single route.
type LoadBalancer struct {
	targets []*Target
	current uint64
	mu      sync.RWMutex
}

func NewLoadBalancer(urls []string) (*LoadBalancer, error) {
	lb := &LoadBalancer{
		targets: make([]*Target, 0),
	}

	for _, u := range urls {
		targetURL, err := url.Parse(u)
		if err != nil {
			return nil, err
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		
		// Custom Director for logging
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			log.Printf("⚖️  LB: Forwarding [%s] %s to %s", req.Method, req.URL.Path, u)
		}

		target := &Target{
			URL:   targetURL,
			Proxy: proxy,
			Alive: true, // Assume alive initially
		}
		lb.targets = append(lb.targets, target)
	}

	// Start health checker
	go lb.healthCheck()

	return lb, nil
}

func (lb *LoadBalancer) healthCheck() {
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ticker.C:
			lb.mu.RLock()
			for _, t := range lb.targets {
				alive := isAlive(t.URL)
				t.SetAlive(alive)
				if !alive {
					log.Printf("🏥 Health Check: Target %s is DOWN", t.URL)
				}
			}
			lb.mu.RUnlock()
		}
	}
}

func isAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	// Filter only healthy targets
	var healthyTargets []*Target
	for _, t := range lb.targets {
		if t.IsAlive() {
			healthyTargets = append(healthyTargets, t)
		}
	}

	if len(healthyTargets) == 0 {
		http.Error(w, "Service Unavailable: No healthy backends found", http.StatusServiceUnavailable)
		return
	}

	// Round Robin logic: increment and wrap
	idx := atomic.AddUint64(&lb.current, 1) % uint64(len(healthyTargets))
	target := healthyTargets[idx]

	target.Proxy.ServeHTTP(w, r)
}
