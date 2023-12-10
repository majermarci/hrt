package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func runTest(test string, ep Endpoint, client *http.Client) {
	// Perform an HTTP request for the endpoint
	req, err := http.NewRequest(ep.Method, ep.URL, strings.NewReader(ep.Body))
	if err != nil {
		fmt.Printf("Failed to create request for endpoint %s: %v\n", test, err)
		os.Exit(1)
	}

	// Add headers to the request
	for key, value := range ep.Headers {
		req.Header.Add(key, value)
	}

	// Send the request using the client
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request for endpoint %s: %v\n", test, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body for endpoint %s: %v", test, err)
		os.Exit(1)
	}

	fmt.Printf("Response for endpoint %s: %v\n", test, resp.Status)
	fmt.Printf("Body: %s\n", body)
}
