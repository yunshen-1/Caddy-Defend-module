package fetchers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OracleFetcher implements the IPRangeFetcher interface for Oracle.
type OracleFetcher struct{}

func (f OracleFetcher) Name() string {
	return "oci"
}

func (f OracleFetcher) Description() string {
	return "Fetches IP ranges for Oracle Cloud Infrastructure services."
}

func (f OracleFetcher) FetchIPRanges() ([]string, error) {
	resp, err := http.Get("https://docs.oracle.com/iaas/tools/public_ip_ranges.json")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Oracle IP ranges: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from Oracle: %v", err)
	}

	var result struct {
		Regions []struct {
			CIDRs []struct {
				CIDR string `json:"cidr"`
			} `json:"cidrs"`
		} `json:"regions"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Oracle JSON: %v", err)
	}

	ipRanges := make([]string, 0, 1000) // default to 1000 IP ranges as an initial capacity
	for _, region := range result.Regions {
		for _, cidr := range region.CIDRs {
			ipRanges = append(ipRanges, cidr.CIDR)
		}
	}

	return ipRanges, nil
}
