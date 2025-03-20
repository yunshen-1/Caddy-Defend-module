package fetchers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// VultrFetcher implements the IPRangeFetcher interface for Vultr.
type VultrFetcher struct{}

func (f VultrFetcher) Name() string {
	return "vultr"
}

func (f VultrFetcher) Description() string {
	return "Fetches IP ranges from Vultr Cloud."
}

func (f VultrFetcher) FetchIPRanges() ([]string, error) {
	const url = "https://geofeed.constant.com/?json"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Vultr IP ranges: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code from Vultr: %d", resp.StatusCode)
	}

	var result struct {
		Description string `json:"description"`
		Email       string `json:"email"`
		Updated     string `json:"updated"`
		Subnets     []struct {
			IPPrefix   string `json:"ip_prefix"`
			Alpha2Code string `json:"alpha2code"`
			Region     string `json:"region"`
			City       string `json:"city"`
			PostalCode string `json:"postal_code"`
		} `json:"subnets"`
		ASN int `json:"asn"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Vultr JSON: %v", err)
	}

	ipRanges := make([]string, 0, len(result.Subnets))
	for _, subnet := range result.Subnets {
		if subnet.IPPrefix != "" {
			ipRanges = append(ipRanges, subnet.IPPrefix)
		}
	}

	return ipRanges, nil
}
