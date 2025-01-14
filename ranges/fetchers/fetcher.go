package fetchers

// IPRangeFetcher defines the interface for fetching IP ranges.
type IPRangeFetcher interface {
	Name() string                     // Returns the name of the service.
	Description() string              // Returns a short description of the service.
	FetchIPRanges() ([]string, error) // Fetches the IP ranges for the service.
}
