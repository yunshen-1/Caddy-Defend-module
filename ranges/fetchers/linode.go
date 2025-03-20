package fetchers

import (
	"bufio"
	"fmt"
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
	const linodeURL = "https://geoip.linode.com/"

	resp, err := http.Get(linodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Linode IP ranges: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code from Linode: %d", resp.StatusCode)
	}

	// Process the response line by line
	var ipRanges []string
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()
		// Skip comment lines that start with #
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Parse the CSV line manually
		fields := strings.Split(line, ",")
		if len(fields) > 0 && fields[0] != "" {
			ipRanges = append(ipRanges, fields[0])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading Linode response: %v", err)
	}

	return ipRanges, nil
}
