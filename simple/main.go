package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"golang.org/x/sync/errgroup"
)

var urls = []string{
	"https://www.codeheim.io",
	"https://golang.org",
	"https://pkg.go.dev/golang.org/x/sync/errgroup",
	"https://invalid-url",
}

func fetchPage(url string, mu *sync.Mutex, responses *map[string]string) error {
	// Fetch the URL content
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("failed to fetch %s: %s\n", url, err)
		return fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	fmt.Printf("Successfully fetched %s\n", url)
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response from %s: %w", url, err)
	}

	// Store the result in the map
	mu.Lock()
	(*responses)[url] = string(body)
	mu.Unlock()

	fmt.Printf("Successfully fetched response body of %s\n", url)

	return nil
}

func main() {
	var g errgroup.Group

	g.SetLimit(2)

	// Create a map to store the responses
	responses := make(map[string]string)
	var mu sync.Mutex

	for _, url := range urls {
		// Start a goroutine for each URL
		g.Go(func() error {
			return fetchPage(url, &mu, &responses)
		})

	}

	// Wait for all goroutines to finish and collect errors
	if err := g.Wait(); err != nil {
		fmt.Println("Error occurred:", err)
	} else {
		fmt.Println("All URLs fetched successfully!")
		// Print the responses
		for url, content := range responses {
			fmt.Printf("Response from %s: %s\n", url, content[:100]) // Print the first 100 characters
		}
	}
}
