package fetchers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// CloudflareFetcher implements the IPRangeFetcher interface for Cloudflare.
type CloudflareFetcher struct{}

func (f CloudflareFetcher) Name() string {
	return "cloudflare"
}

func (f CloudflareFetcher) Description() string {
	return "Fetches IP ranges used by Cloudflare services."
}

func (f CloudflareFetcher) FetchIPRanges() ([]string, error) {
	const url = "https://api.cloudflare.com/client/v4/ips"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Cloudflare IP ranges: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code from Cloudflare: %d", resp.StatusCode)
	}

	var result struct {
		Result struct {
			IPv4CIDRs []string `json:"ipv4_cidrs"`
			IPv6CIDRs []string `json:"ipv6_cidrs"`
		} `json:"result"`
		Success bool `json:"success"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Cloudflare JSON: %v", err)
	}

	// Combine IPv4 and IPv6 ranges
	ipRanges := make([]string, 0, len(result.Result.IPv4CIDRs)+len(result.Result.IPv6CIDRs))
	ipRanges = append(ipRanges, result.Result.IPv4CIDRs...)
	ipRanges = append(ipRanges, result.Result.IPv6CIDRs...)

	return ipRanges, nil
}
