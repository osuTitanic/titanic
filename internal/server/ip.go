package server

import (
	"net"
	"net/http"
	"strings"
)

// Cloudflare specific header for original client IP
var CloudflareIPHeader = "CF-Connecting-IP"

// Common reverse proxy headers
var ipHeaders = []string{
	"X-Forwarded-For",
	"X-Real-IP",
	"Forwarded",
}

// GetRequestIP tries to extract the real client IP address from the request headers
func GetRequestIP(r *http.Request) string {
	if ip := r.Header.Get(CloudflareIPHeader); ip != "" {
		return ip
	}

	for _, h := range ipHeaders {
		if val := r.Header.Get(h); val != "" {
			// X-Forwarded-For may be a comma-separated list -> take the first entry
			ip := strings.TrimSpace(strings.Split(val, ",")[0])
			if ip != "" {
				return ip
			}
		}
	}

	// Fallback to remote address
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
