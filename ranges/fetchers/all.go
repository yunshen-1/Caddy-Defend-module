package fetchers

// AllFetcher implements the IPRangeFetcher interface for all network ranges.
type AllFetcher struct{}

func (f AllFetcher) Name() string {
	return "All"
}
func (f AllFetcher) Description() string {
	return "Every IP address in existence."
}
func (f AllFetcher) FetchIPRanges() ([]string, error) {
	return []string{
		"::/0",
		"0.0.0.0/0",
	}, nil
}
