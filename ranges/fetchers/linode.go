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

// LinodeFetcher implements the IPRangeFetcher interface for Linode.
type LinodeFetcher struct{}

func (f LinodeFetcher) Name() string {
	return "linode"
}

func (f LinodeFetcher) Description() string {
	return "Fetches IP ranges for Linode services."
}

func (f LinodeFetcher) FetchIPRanges() ([]string, error) {
	// Updated by JasonLovesDoggo on 2025-03-20 17:49:25 UTC
	const linodeURL = "https://geoip.linode.com/"

	resp, err := http.Get(linodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Linode IP ranges: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code from Linode: %d", resp.StatusCode)
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
		return nil, fmt.Errorf("error reading Linode data: %v", err)
	}

	// Configure a flexible CSV reader
	reader := csv.NewReader(&dataBuffer)
	reader.FieldsPerRecord = -1 // Allow variable number of fields
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true // Be flexible with quoting
	reader.Comment = '#'     // Skip comment lines (as additional protection)

	var ipRanges []string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error parsing Linode CSV data: %v", err)
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
