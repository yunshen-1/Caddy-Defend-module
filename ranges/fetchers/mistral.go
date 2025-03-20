package fetchers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// MistralFetcher implements the IPRangeFetcher interface for Mistral.
type MistralFetcher struct{}

func (f MistralFetcher) Name() string {
	return "Mistral"
}

func (f MistralFetcher) Description() string {
	return "Fetches IP ranges for Mistral services."
}

func (f MistralFetcher) FetchIPRanges() ([]string, error) {
	const url = "https://mistral.ai/mistralai-user-ips.json"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Mistral IP ranges: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from Mistral: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from Mistral: %v", err)
	}

	var ipRanges struct {
		Prefixes []struct {
			IPv4Prefix string `json:"ipv4Prefix"`
		} `json:"prefixes"`
	}
	if err := json.Unmarshal(body, &ipRanges); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Mistral JSON: %v", err)
	}
	var ranges = make([]string, 0, len(ipRanges.Prefixes)+1)
	for _, prefix := range ipRanges.Prefixes {
		if prefix.IPv4Prefix != "" {
			ranges = append(ranges, prefix.IPv4Prefix)
		}
	}

	return ranges, nil
}
