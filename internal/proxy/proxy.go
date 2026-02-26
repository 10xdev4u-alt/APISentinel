package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

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
