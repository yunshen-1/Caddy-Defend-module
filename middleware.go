package caddydefender

import (
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"net"
	"net/http"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/jasonlovesdoggo/caddy-defender/ranges/data"
)

// ServeHTTP implements the middleware logic.
func (m DefenderMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	clientIP := net.ParseIP(r.RemoteAddr)
	if clientIP == nil {
		return caddyhttp.Error(http.StatusForbidden, fmt.Errorf("invalid client IP"))
	}

	// Check if the client IP is in any of the embedded ranges
	for _, ranges := range data.IPRanges {
		for _, cidr := range ranges {
			_, ipNet, err := net.ParseCIDR(cidr)
			if err != nil {
				caddy.Log().Error(fmt.Sprintf("Invalid CIDR: %v", err))
				continue
			}
			if ipNet.Contains(clientIP) {
				return m.Responder.Respond(w, r)
			}
		}
	}

	// Check if the client IP is in any of the additional ranges
	for _, cidr := range m.AdditionalRanges {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			caddy.Log().Error(fmt.Sprintf("Invalid CIDR: %v", err))
			continue
		}
		if ipNet.Contains(clientIP) {
			return m.Responder.Respond(w, r)
		}
	}

	// IP is not in any of the ranges, proceed to the next handler
	return next.ServeHTTP(w, r)
}
