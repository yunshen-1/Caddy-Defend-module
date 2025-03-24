package aws

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// IPRanges represents the structure of the AWS IP ranges JSON file.
type IPRanges struct {
	SyncToken  string `json:"syncToken"`
	CreateDate string `json:"createDate"`
	Prefixes   []struct {
		IPPrefix string `json:"ip_prefix"`
		Region   string `json:"region"`
		Service  string `json:"service"`
	} `json:"prefixes"`
	IPv6Prefixes []struct {
		IPv6Prefix string `json:"ipv6_prefix"`
		Region     string `json:"region"`
		Service    string `json:"service"`
	} `json:"ipv6_prefixes"`
}

// fetchAWSIPRanges fetches and parses the AWS IP ranges JSON file.
func fetchAWSIPRanges(region, service string) ([]string, error) {
	const url = "https://ip-ranges.amazonaws.com/ip-ranges.json"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch AWS IP ranges from %s: %v", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s: %v", url, err)
	}

	var ipRanges IPRanges
	if err := json.Unmarshal(body, &ipRanges); err != nil {
		return nil, fmt.Errorf("failed to unmarshal AWS IP ranges JSON: %v", err)
	}

	// Extract IP ranges for the specified region and service
	var ranges []string
	for _, prefix := range ipRanges.Prefixes {
		if (region == "" || prefix.Region == region) && (service == "" || prefix.Service == service) {
			ranges = append(ranges, prefix.IPPrefix)
		}
	}
	for _, prefix := range ipRanges.IPv6Prefixes {
		if (region == "" || prefix.Region == region) && (service == "" || prefix.Service == service) {
			ranges = append(ranges, prefix.IPv6Prefix)
		}
	}

	return ranges, nil
}
