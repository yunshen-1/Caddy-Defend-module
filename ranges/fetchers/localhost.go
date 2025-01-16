package fetchers

// LocalhostFetcher implements the IPRangeFetcher interface for Localhost.
type LocalhostFetcher struct{}

func (f LocalhostFetcher) Name() string {
	return "Localhost"
}
func (f LocalhostFetcher) Description() string {
	return "Hardcoded IP ranges for Localhost. Used in development."
}
func (f LocalhostFetcher) FetchIPRanges() ([]string, error) {

	return []string{
		"127.0.0.0/8", // IPv4 localhost range
		"::1/128",     // IPv6 localhost range
	}, nil
}
