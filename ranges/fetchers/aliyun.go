package fetchers

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

// AliyunFetcher implements the IPRangeFetcher interface for Alibaba Cloud (Aliyun).
type AliyunFetcher struct{}

func (f AliyunFetcher) Name() string {
	return "Aliyun"
}

func (f AliyunFetcher) Description() string {
	return "Fetches IP ranges for Alibaba Cloud (Aliyun) services."
}

func (f AliyunFetcher) FetchIPRanges() ([]string, error) {
	const aliyunURL = "https://raw.githubusercontent.com/sakib-m/IP-Prefix-List/main/ALIBABA/only_ip_blocks.txt"

	resp, err := http.Get(aliyunURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Alibaba Cloud IP ranges: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code from Alibaba Cloud IP list: %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	var ipRanges []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		ipRanges = append(ipRanges, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading Alibaba Cloud IP ranges: %v", err)
	}

	return ipRanges, nil
}
