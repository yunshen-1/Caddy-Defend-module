package utils

import (
	"errors"
	"fmt"
	"github.com/gaissmai/bart"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/jasonlovesdoggo/caddy-defender/ranges/data"
	"go.uber.org/zap"
	"net"
	"net/netip"
	"time"
)

const MaxKeys = 10000

// IPChecker holds the CIDR ranges in an optimized structure for fast lookups
type IPChecker struct {
	table *bart.Table[struct{}] // Using empty struct as value since we only need existence check
	cache *expirable.LRU[string, bool]
	log   *zap.Logger
}

// NewIPChecker creates a new IPChecker instance with preprocessed CIDR ranges
func NewIPChecker(cidrRanges []string, log *zap.Logger) *IPChecker {
	return &IPChecker{
		table: buildTable(cidrRanges, log),
		cache: expirable.NewLRU[string, bool](MaxKeys, nil, time.Minute*10),
		log:   log,
	}
}

// buildTable initializes the radix tree with CIDR ranges during provisioning
func buildTable(cidrRanges []string, log *zap.Logger) *bart.Table[struct{}] {
	table := &bart.Table[struct{}]{}
	for _, cidr := range cidrRanges {
		// Handle predefined range groups
		if ranges, ok := data.IPRanges[cidr]; ok {
			for _, predefinedCIDR := range ranges {
				if err := insertCIDR(table, predefinedCIDR, log); err != nil {
					log.Error("invalid predefined CIDR",
						zap.String("group", cidr),
						zap.String("cidr", predefinedCIDR),
						zap.Error(err))
				}
			}
			continue
		}

		// Handle direct CIDR specifications
		if err := insertCIDR(table, cidr, log); err != nil {
			log.Error("invalid CIDR specification",
				zap.String("cidr", cidr),
				zap.Error(err))
		}
	}
	return table
}

// insertCIDR safely parses and inserts a CIDR into the radix tree
func insertCIDR(table *bart.Table[struct{}], cidr string, log *zap.Logger) error {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return fmt.Errorf("failed to parse CIDR: %w", err)
	}
	table.Insert(prefix, struct{}{})
	return nil
}

// IPInRanges checks if an IP address matches any CIDR range with caching
func (c *IPChecker) IPInRanges(clientIP net.IP) bool {
	normalizedIP := normalizeIP(clientIP)

	// Cache check
	if val, ok := c.cache.Get(normalizedIP); ok {
		return val
	}

	// Convert to netip.Addr for radix tree lookup
	ipAddr, err := netIPToNetipAddr(clientIP)
	if err != nil {
		c.log.Error("invalid IP address",
			zap.String("ip", clientIP.String()),
			zap.Error(err))
		return false
	}

	// Radix tree lookup
	result := c.table.Contains(ipAddr)

	// Update cache
	c.cache.Add(normalizedIP, result)
	return result
}

// normalizeIP ensures consistent string representation for caching
func normalizeIP(ip net.IP) string {
	if ip4 := ip.To4(); ip4 != nil {
		return ip4.String()
	}
	return ip.To16().String()
}

// netIPToNetipAddr converts net.IP to netip.Addr with validation
func netIPToNetipAddr(ip net.IP) (netip.Addr, error) {
	addr, ok := netip.AddrFromSlice(ip)
	if !ok {
		return netip.Addr{}, errors.New("invalid IP address format")
	}
	return addr, nil
}
