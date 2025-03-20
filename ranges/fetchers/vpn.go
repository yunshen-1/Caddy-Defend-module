package fetchers

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

// VPNFetcher implements the IPRangeFetcher interface for known VPN services.
type VPNFetcher struct{}

func (f VPNFetcher) Name() string {
	return "vpn"
}

func (f VPNFetcher) Description() string {
	return "Fetches IP ranges of known VPN services."
}

func (f VPNFetcher) FetchIPRanges() ([]string, error) {
	const vpnURL = "https://cdn.jsdelivr.net/gh/X4BNet/lists_vpn@main/output/vpn/ipv4.txt"

	resp, err := http.Get(vpnURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch VPN IP ranges: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code from VPN list: %d", resp.StatusCode)
	}

	var ipRanges []string
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			ipRanges = append(ipRanges, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading VPN IP ranges: %v", err)
	}

	return ipRanges, nil
}
