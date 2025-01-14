package fetchers

// DeepSeekFetcher implements the IPRangeFetcher interface for DeepSeek.
type DeepSeekFetcher struct{}

func (f DeepSeekFetcher) Name() string {
	return "DeepSeek"
}
func (f DeepSeekFetcher) Description() string {
	return "Hardcoded IP ranges for DeepSeek services."
}
func (f DeepSeekFetcher) FetchIPRanges() ([]string, error) {
	// https://discuss.deepsource.com/t/incoming-adding-new-ip-addresses-to-deepsources-ip-range/667

	return []string{
		"35.225.112.198/32",
		"34.42.70.44/32",
		"104.154.172.152/32",
	}, nil
}
