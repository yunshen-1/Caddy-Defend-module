package main

import (
	"encoding/json"
	"fmt"
	"github.com/jasonlovesdoggo/caddy-defender/ranges/fetchers"
	"log"
	"os"
	"sync"
)

func main() {
	// Create an array of all IP range fetchers
	fetchersList := []fetchers.IPRangeFetcher{
		fetchers.OpenAIFetcher{},
		fetchers.DeepSeekFetcher{},
		fetchers.GithubCopilotFetcher{},
	}

	// Create a map to hold the IP ranges
	ipRanges := make(map[string][]string)

	// Use a WaitGroup to wait for all fetchers to complete
	var wg sync.WaitGroup
	wg.Add(len(fetchersList))

	// Use a mutex to safely update the ipRanges map
	var mu sync.Mutex

	// Start fetching IP ranges concurrently
	for _, fetcher := range fetchersList {
		go func(f fetchers.IPRangeFetcher) {
			defer wg.Done()

			// Print the start of the fetching process
			fmt.Printf("🚀 Starting %s: %s\n", f.Name(), f.Description())

			// Fetch the IP ranges
			ranges, err := f.FetchIPRanges()
			if err != nil {
				fmt.Printf("❌ Error fetching %s: %v\n", f.Name(), err)
				return
			}

			// Update the map with the fetched ranges
			mu.Lock()
			ipRanges[f.Name()] = ranges
			mu.Unlock()

			// Print the completion of the fetching process
			fmt.Printf("✅ Completed %s: Fetched %d IP ranges\n", f.Name(), len(ranges))
		}(fetcher)
	}

	wg.Wait()

	// Convert the map to JSON
	jsonData, err := json.MarshalIndent(ipRanges, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal IP ranges to JSON: %v", err)
	}

	// Write the JSON data to a file
	err = os.WriteFile("output.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON to file: %v", err)
	}

	fmt.Println("\n🎉 All IP ranges have been successfully written to output.json")
}
