package aws

import "fmt"

// RegionFetcher AWSRegionFetcher implements the IPRangeFetcher interface for AWS regions.
type RegionFetcher struct {
	Region string // The AWS region to fetch IP ranges for
}

func (f RegionFetcher) Name() string {
	return fmt.Sprintf("AWS-%s", f.Region)
}

func (f RegionFetcher) Description() string {
	return fmt.Sprintf("Fetches IP ranges for AWS services in the %s region.", f.Region)
}

func (f RegionFetcher) FetchIPRanges() ([]string, error) {
	// Fetch AWS IP ranges for the specified region
	return fetchAWSIPRanges(f.Region, "")
}
