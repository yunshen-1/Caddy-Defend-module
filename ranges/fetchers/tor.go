package fetchers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

// TorFetcher implements the IPRangeFetcher interface for Tor exit nodes.
type TorFetcher struct{}

func (f TorFetcher) Name() string {
	return "tor"
}

func (f TorFetcher) Description() string {
	return "Fetches IP addresses of Tor exit nodes."
}

func (f TorFetcher) FetchIPRanges() ([]string, error) {
	const torURL = "https://cdn.jsdelivr.net/gh/alireza-rezaee/tor-nodes@main/latest.exits.csv"

	resp, err := http.Get(torURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Tor exit nodes: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code from Tor node list: %d", resp.StatusCode)
	}

	// Parse the CSV file
	r := csv.NewReader(resp.Body)
	// Skip the header row
	if _, err := r.Read(); err != nil {
		return nil, fmt.Errorf("failed to read header from Tor exit nodes CSV: %v", err)
	}

	var ipRanges = make([]string, resp.ContentLength/4) // 4 bytes per IP address (assuming ipv4)

	// Read the rest of the records
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading Tor exit nodes CSV: %v", err)
		}

		// The IP address is in the second column (index 1)
		if len(record) > 1 {
			ipStr := strings.TrimSpace(record[1])
			if ipStr == "" {
				continue
			}

			// Convert IP address to CIDR notation
			ip := net.ParseIP(ipStr)
			if ip == nil {
				// Skip invalid IPs
				continue
			}

			var cidr string
			if strings.Contains(ipStr, ":") {
				// IPv6 address - use /128
				cidr = fmt.Sprintf("%s/128", ipStr)
			} else {
				// IPv4 address - use /32
				cidr = fmt.Sprintf("%s/32", ipStr)
			}

			ipRanges = append(ipRanges, cidr)
		}
	}

	return ipRanges, nil
}
