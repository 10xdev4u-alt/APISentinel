package logger

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// IPLocation represents the geographical information for an IP address.
type IPLocation struct {
	Country string `json:"country"`
	City    string `json:"city"`
	ISP     string `json:"isp"`
}

var (
	geoCache = make(map[string]*IPLocation)
	geoMu    sync.RWMutex
)

// GetLocation looks up the location for an IP address using ip-api.com.
func GetLocation(ip string) (*IPLocation, error) {
	// Strip port if present
	host, _, err := net.SplitHostPort(ip)
	if err == nil {
		ip = host
	}

	// 1. Check Cache
	geoMu.RLock()
	if loc, ok := geoCache[ip]; ok {
		geoMu.RUnlock()
		return loc, nil
	}
	geoMu.RUnlock()

	// 2. Lookup (Skip if local IP)
	if ip == "::1" || ip == "127.0.0.1" || ip == "localhost" || strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") {
		return &IPLocation{Country: "Local", City: "Dev Machine", ISP: "Internal"}, nil
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://ip-api.com/json/%s", ip))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var loc IPLocation
	if err := json.NewDecoder(resp.Body).Decode(&loc); err != nil {
		return nil, err
	}

	// 3. Save to Cache
	geoMu.Lock()
	geoCache[ip] = &loc
	geoMu.Unlock()

	return &loc, nil
}
