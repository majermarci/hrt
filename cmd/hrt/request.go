package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	RequestName  string
	StatusCode   string
	ResponseBody string
	Method       string
	URL          string
}

func runTest(test string, ep Endpoint, client *http.Client) Response {
	if ep.Method == "" {
		ep.Method = "GET"
	}

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

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body for endpoint %s: %v", test, err)
		os.Exit(1)
	}

	// Try to unmarshal the body into a JSON object
	var jsonObj map[string]interface{}
	if err := json.Unmarshal(body, &jsonObj); err == nil {
		// If successful, re-marshal it with indentation
		body, _ = json.MarshalIndent(jsonObj, "", "  ")
	}

	// fmt.Printf("Status Code for '%s': %v\n", test, resp.Status)
	// fmt.Printf("Response Body: \n%s\n", body)

	// Return the response
	return Response{
		RequestName:  test,
		Method:       ep.Method,
		StatusCode:   resp.Status,
		ResponseBody: string(body),
		URL:          ep.URL,
	}
}
