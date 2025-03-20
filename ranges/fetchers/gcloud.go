package fetchers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GCloudFetcher implements the IPRangeFetcher interface for GCP IP ranges.
type GCloudFetcher struct{}

func (f GCloudFetcher) Name() string {
	return "GCloud"
}

func (f GCloudFetcher) Description() string {
	return "Fetches IP ranges for Google Cloud Platform (GCP) services."
}

func (f GCloudFetcher) FetchIPRanges() ([]string, error) {
	// Fetch all GCP IP ranges
	return fetchGCloudIPRanges()
}

// GCloudIPRanges represents the structure of the GCP IP ranges JSON file.
type GCloudIPRanges struct {
	SyncToken    string `json:"syncToken"`
	CreationTime string `json:"creationTime"`
	Prefixes     []struct {
		IPv4Prefix string `json:"ipv4Prefix"`
		IPv6Prefix string `json:"ipv6Prefix"`
		Service    string `json:"service"`
		Scope      string `json:"scope"`
	} `json:"prefixes"`
}

// fetchGCloudIPRanges fetches and parses the GCP IP ranges JSON file.
func fetchGCloudIPRanges() ([]string, error) {
	const url = "https://www.gstatic.com/ipranges/cloud.json"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GCP IP ranges from %s: %v", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s: %v", url, err)
	}

	var ipRanges GCloudIPRanges
	if err := json.Unmarshal(body, &ipRanges); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GCP IP ranges JSON: %v", err)
	}

	// Extract all IP ranges (both IPv4 and IPv6)
	var ranges []string
	for _, prefix := range ipRanges.Prefixes {
		if prefix.IPv4Prefix != "" {
			ranges = append(ranges, prefix.IPv4Prefix)
		}
		if prefix.IPv6Prefix != "" {
			ranges = append(ranges, prefix.IPv6Prefix)
		}
	}

	return ranges, nil
}
