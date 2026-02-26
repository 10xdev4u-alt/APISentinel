package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"
	"strings"
)

// MultiTargetProxy routes requests to different LoadBalancers based on path prefixes.
type MultiTargetProxy struct {
	Routes map[string]*LoadBalancer
}

func NewMultiTargetProxy() *MultiTargetProxy {
	return &MultiTargetProxy{
		Routes: make(map[string]*LoadBalancer),
	}
}

func (m *MultiTargetProxy) AddRoute(prefix string, targets []string) error {
	lb, err := NewLoadBalancer(targets)
	if err != nil {
		return err
	}

	m.Routes[prefix] = lb
	return nil
}

func (m *MultiTargetProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get prefixes and sort them by length descending
	prefixes := make([]string, 0, len(m.Routes))
	for k := range m.Routes {
		prefixes = append(prefixes, k)
	}
	sort.Slice(prefixes, func(i, j int) bool {
		return len(prefixes[i]) > len(prefixes[j])
	})

	for _, prefix := range prefixes {
		if strings.HasPrefix(r.URL.Path, prefix) {
			m.Routes[prefix].ServeHTTP(w, r)
			return
		}
	}

	// Default fallback if no route matches (could be handled in main)
	http.Error(w, "Not Found: No route matches this path", http.StatusNotFound)
}

// NewProxy creates a new ReverseProxy for the given target URL.
func NewProxy(target string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	// Custom Director: This is where we can modify the request before it's sent.
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		log.Printf("🛡️ API Sentinel: Forwarding request [%s] %s to %s", req.Method, req.URL.Path, target)
	}

	return proxy, nil
}
