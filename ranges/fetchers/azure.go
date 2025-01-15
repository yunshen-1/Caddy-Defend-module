package fetchers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// AzurePublicCloudFetcher implements the IPRangeFetcher interface for Azure Public Cloud.
type AzurePublicCloudFetcher struct{}

func (f AzurePublicCloudFetcher) Name() string {
	return "AzurePublicCloud"
}

func (f AzurePublicCloudFetcher) Description() string {
	return "Fetches IP ranges for Azure Public Cloud services."
}

func (f AzurePublicCloudFetcher) FetchIPRanges() ([]string, error) {
	// https://www.microsoft.com/en-us/download/details.aspx?id=56519
	const downloadURL = "https://download.microsoft.com/download/7/1/D/71D86715-5596-4529-9B13-DA13A5DE5B63/ServiceTags_Public_20250113.json"
	resp, err := http.Get(downloadURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Azure Public Cloud IP ranges: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from Azure Public Cloud: %v", err)
	}

	type AzureIPRanges struct {
		Values []struct {
			Name       string `json:"name"`
			Properties struct {
				Platform        string   `json:"platform"`
				SystemService   string   `json:"systemService"`
				AddressPrefixes []string `json:"addressPrefixes"`
			} `json:"properties"`
		} `json:"values"`
	}

	var ipRanges AzureIPRanges
	if err := json.Unmarshal(body, &ipRanges); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Azure Public Cloud JSON: %v", err)
	}

	// Filter out the "Public" cloud IPs
	publicIPs := []string{}
	for _, value := range ipRanges.Values {
		if (strings.EqualFold(value.Properties.Platform, "Azure")) && (strings.EqualFold(value.Properties.SystemService, "ActionGroup")) {
			publicIPs = append(publicIPs, value.Properties.AddressPrefixes...)
		}
	}

	return publicIPs, nil
}
