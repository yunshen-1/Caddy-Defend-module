package fetchers

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// AliyunFetcher implements the IPRangeFetcher interface for Alibaba Cloud (Aliyun).
type AliyunFetcher struct{}

func (f AliyunFetcher) Name() string {
	return "aliyun"
}

func (f AliyunFetcher) Description() string {
	return "Fetches IP ranges for Alibaba Cloud (Aliyun) services."
}

func (f AliyunFetcher) FetchIPRanges() ([]string, error) {
	const aliyunURL = "https://cdn.jsdelivr.net/gh/sakib-m/IP-Prefix-List@main/ALIBABA/only_ip_blocks.txt"

	resp, err := http.Get(aliyunURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Alibaba Cloud IP ranges: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code from Alibaba Cloud IP list: %d", resp.StatusCode)
	}

	// Pre-process the data to remove comment lines
	var dataBuffer bytes.Buffer
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			dataBuffer.WriteString(line + "\n")
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading Alibaba Cloud data: %v", err)
	}

	// Configure a flexible CSV reader
	reader := csv.NewReader(&dataBuffer)
	reader.FieldsPerRecord = -1 // Allow variable number of fields
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true // Be flexible with quoting
	reader.Comment = '#'     // Skip comment lines (as additional protection)

	// Preallocate the slice with an estimated size
	ipRanges := make([]string, 0, 1700)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error parsing Alibaba Cloud CSV data: %v", err)
		}

		// Extract the IP range from the first field if available
		if len(record) > 0 && record[0] != "" {
			ipRange := strings.TrimSpace(record[0])
			if ipRange != "" {
				ipRanges = append(ipRanges, ipRange)
			}
		}
	}

	return ipRanges, nil
}
