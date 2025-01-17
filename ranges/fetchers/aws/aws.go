package aws

// AWSFetcher implements the IPRangeFetcher interface for AWS global IP ranges.
type AWSFetcher struct{}

func (f AWSFetcher) Name() string {
	return "AWS"
}

func (f AWSFetcher) Description() string {
	return "Fetches global IP ranges for AWS services."
}

func (f AWSFetcher) FetchIPRanges() ([]string, error) {
	// Fetch all AWS IP ranges (no region or service filter)
	return fetchAWSIPRanges("", "")
}
