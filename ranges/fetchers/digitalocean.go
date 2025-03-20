package fetchers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
)

// DigitalOceanFetcher implements the IPRangeFetcher interface for Digital Ocean.
type DigitalOceanFetcher struct{}

func (f DigitalOceanFetcher) Name() string {
	return "digitalocean"
}

func (f DigitalOceanFetcher) Description() string {
	return "Fetches IP ranges for Digital Ocean services."
}

func (f DigitalOceanFetcher) FetchIPRanges() ([]string, error) {
	const doURL = "https://digitalocean.com/geo/google.csv"

	resp, err := http.Get(doURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Digital Ocean IP ranges: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code from Digital Ocean: %d", resp.StatusCode)
	}

	// Parse the CSV file
	csvReader := csv.NewReader(resp.Body)

	// No official headers in the file, so skip first row handling
	var ipRanges []string
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading Digital Ocean CSV: %v", err)
		}

		// The first field is the IP prefix
		if len(record) > 0 && record[0] != "" {
			ipRanges = append(ipRanges, record[0])
		}
	}

	return ipRanges, nil
}
