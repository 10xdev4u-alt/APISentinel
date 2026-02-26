package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

// Target represents a single backend server.
type Target struct {
	URL   *url.URL
	Proxy *httputil.ReverseProxy
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

		lb.targets = append(lb.targets, &Target{
			URL:   targetURL,
			Proxy: proxy,
		})
	}

	return lb, nil
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	if len(lb.targets) == 0 {
		http.Error(w, "Service Unavailable: No backends found", http.StatusServiceUnavailable)
		return
	}

	// Round Robin logic: increment and wrap
	idx := atomic.AddUint64(&lb.current, 1) % uint64(len(lb.targets))
	target := lb.targets[idx]

	target.Proxy.ServeHTTP(w, r)
}
