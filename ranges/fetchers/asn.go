package fetchers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ASNFetcher implements the IPRangeFetcher interface for specific ASNs.
type ASNFetcher struct {
	ASNs []string // List of ASNs in AS#### format
}

func (f ASNFetcher) Name() string {
	return "asn"
}

func (f ASNFetcher) Description() string {
	return "Fetches IP ranges for specific Autonomous System Numbers (ASNs)."
}

func (f ASNFetcher) FetchIPRanges() ([]string, error) {
	if len(f.ASNs) == 0 {
		return nil, fmt.Errorf("no ASNs provided to fetch")
	}

	ipRanges := make([]string, 0)

	for _, asn := range f.ASNs {
		url := fmt.Sprintf("https://api.hackertarget.com/aslookup/?q=%s", asn)
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch IP ranges for ASN %s: %v", asn, err)
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("received non-200 status code %d for ASN %s", resp.StatusCode, asn)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body for ASN %s: %v", asn, err)
		}

		// Split the response by newlines
		lines := strings.Split(string(body), "\n")

		// Skip the first line as it contains the ASN info rather than IP ranges
		if len(lines) > 1 {
			for i := 1; i < len(lines); i++ {
				ipRange := strings.TrimSpace(lines[i])
				if ipRange != "" {
					ipRanges = append(ipRanges, ipRange)
				}
			}
		}
	}

	return ipRanges, nil
}

// NewASNFetcher creates a new ASNFetcher with the specified ASNs.
func NewASNFetcher(asns []string) *ASNFetcher {
	return &ASNFetcher{
		ASNs: asns,
	}
}
