package caddydefender

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"net/http"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/jasonlovesdoggo/caddy-defender/ranges/data"
)

// ServeHTTP implements the middleware logic.
func (m DefenderMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Split the RemoteAddr into IP and port
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		m.log.Error("Invalid client IP format", zap.String("ip", r.RemoteAddr))
		return caddyhttp.Error(http.StatusForbidden, fmt.Errorf("invalid client IP format"))
	}

	clientIP := net.ParseIP(host)
	m.log.Debug("client IP", zap.String("ip", clientIP.String()))
	if clientIP == nil {
		m.log.Error("Invalid client IP", zap.String("ip", host))
		return caddyhttp.Error(http.StatusForbidden, fmt.Errorf("invalid client IP"))
	}

	// Check if the client IP is in any of the additional ranges
	for _, cidr := range m.AdditionalRanges {
		// If the range is a predefined key (e.g., "openai"), use the corresponding CIDRs
		if ranges, ok := data.IPRanges[cidr]; ok {
			for _, predefinedCIDR := range ranges {
				_, ipNet, err := net.ParseCIDR(predefinedCIDR)
				if err != nil {
					m.log.Error(fmt.Sprintf("Invalid predefined CIDR: %v", err))
					continue
				}
				if ipNet.Contains(clientIP) {
					return m.responder.Respond(w, r)
				}
			}
		} else {
			// Otherwise, treat it as a custom CIDR
			_, ipNet, err := net.ParseCIDR(cidr)
			if err != nil {
				m.log.Error(fmt.Sprintf("Invalid CIDR: %v", err))
				continue
			}
			if ipNet.Contains(clientIP) {
				return m.responder.Respond(w, r)
			}
		}
	}

	// IP is not in any of the ranges, proceed to the next handler
	return next.ServeHTTP(w, r)
}
