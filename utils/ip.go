package utils

import (
	"fmt"
	"github.com/jasonlovesdoggo/caddy-defender/ranges/data"
	"go.uber.org/zap"
	"net"
)

// IPInRanges checks if the given IP is within any of the provided CIDR ranges.
// It returns true if the IP is in any of the ranges, false otherwise.
func IPInRanges(clientIP net.IP, cidrRanges []string, log *zap.Logger) bool {
	for _, cidr := range cidrRanges {
		// If the range is a predefined key (e.g., "openai"), use the corresponding CIDRs
		if ranges, ok := data.IPRanges[cidr]; ok {
			for _, predefinedCIDR := range ranges {
				_, ipNet, err := net.ParseCIDR(predefinedCIDR)
				if err != nil {
					log.Error(fmt.Sprintf("Invalid predefined CIDR: %v", err))
					continue
				}
				if ipNet.Contains(clientIP) {
					return true
				}
			}
		} else {
			// Otherwise, treat it as a custom CIDR
			_, ipNet, err := net.ParseCIDR(cidr)
			if err != nil {
				log.Error(fmt.Sprintf("Invalid CIDR: %v", err))
				continue
			}
			if ipNet.Contains(clientIP) {
				return true
			}
		}
	}
	return false
}
