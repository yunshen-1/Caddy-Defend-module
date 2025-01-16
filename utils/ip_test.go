package utils

import (
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/jasonlovesdoggo/caddy-defender/ranges/data"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Test data
var (
	validCIDRs = []string{
		"192.168.1.0/24",
		"10.0.0.0/8",
		"2001:db8::/32",
	}
	invalidCIDRs = []string{
		"invalid-cidr",
		"192.168.1.0/33", // Invalid subnet mask
	}
	predefinedCIDRs = map[string][]string{
		"openai": {
			"203.0.113.0/24",
			"2001:db8:1::/48",
		},
	}
)

// Mock logger for testing
var testLogger = zap.NewNop()

// TestRawIPInRanges tests the rawIPInRanges function.
func TestRawIPInRanges(t *testing.T) {
	tests := []struct {
		name       string
		ip         string
		cidrRanges []string
		expected   bool
	}{
		{
			name:       "IPv4 in range",
			ip:         "192.168.1.100",
			cidrRanges: validCIDRs,
			expected:   true,
		},
		{
			name:       "IPv4 not in range",
			ip:         "192.168.2.100",
			cidrRanges: validCIDRs,
			expected:   false,
		},
		{
			name:       "IPv6 in range",
			ip:         "2001:db8::1",
			cidrRanges: validCIDRs,
			expected:   true,
		},
		{
			name:       "IPv6 not in range",
			ip:         "2001:db8:1::1",
			cidrRanges: []string{"2001:db8::/48"}, // Narrower range
			expected:   false,
		},
		{
			name:       "Invalid CIDR",
			ip:         "192.168.1.100",
			cidrRanges: invalidCIDRs,
			expected:   false,
		},
		{
			name:       "Predefined CIDR (IPv4)",
			ip:         "203.0.113.10",
			cidrRanges: []string{"openai"},
			expected:   true,
		},
		{
			name:       "Predefined CIDR (IPv6)",
			ip:         "2001:db8:1::10",
			cidrRanges: []string{"openai"},
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientIP := net.ParseIP(tt.ip)
			assert.NotNil(t, clientIP, "Failed to parse IP")

			// Mock predefined CIDRs
			data.IPRanges = predefinedCIDRs

			result := rawIPInRanges(clientIP, tt.cidrRanges, testLogger)
			assert.Equal(t, tt.expected, result, "Unexpected result for IP %s", tt.ip)
		})
	}
}

// TestIPInRanges tests the IPInRanges function, including caching behavior.
func TestIPInRanges(t *testing.T) {
	tests := []struct {
		name       string
		ip         string
		cidrRanges []string
		expected   bool
	}{
		{
			name:       "IPv4 in range (cached)",
			ip:         "192.168.1.100",
			cidrRanges: validCIDRs,
			expected:   true,
		},
		{
			name:       "IPv4 not in range (cached)",
			ip:         "192.168.2.100",
			cidrRanges: validCIDRs,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientIP := net.ParseIP(tt.ip)
			assert.NotNil(t, clientIP, "Failed to parse IP")

			// Mock predefined CIDRs
			data.IPRanges = predefinedCIDRs

			// First call (not cached)
			result := IPInRanges(clientIP, tt.cidrRanges, testLogger)
			assert.Equal(t, tt.expected, result, "Unexpected result for IP %s (first call)", tt.ip)

			// Second call (cached)
			result = IPInRanges(clientIP, tt.cidrRanges, testLogger)
			assert.Equal(t, tt.expected, result, "Unexpected result for IP %s (second call)", tt.ip)
		})
	}
}

// TestIPInRangesCacheExpiration tests the cache expiration behavior.
func TestIPInRangesCacheExpiration(t *testing.T) {
	// Set a short cache expiration time for testing
	cache = expirable.NewLRU[string, bool](MaxKeys, nil, time.Millisecond*10)

	clientIP := net.ParseIP("192.168.1.100")
	assert.NotNil(t, clientIP, "Failed to parse IP")

	// Mock predefined CIDRs
	data.IPRanges = predefinedCIDRs

	// First call (not cached)
	result := IPInRanges(clientIP, validCIDRs, testLogger)
	assert.True(t, result, "Expected IP to be in range (first call)")

	// Wait for cache to expire
	time.Sleep(time.Millisecond * 20)

	// Second call (cache expired)
	result = IPInRanges(clientIP, validCIDRs, testLogger)
	assert.True(t, result, "Expected IP to be in range (second call, cache expired)")
}
