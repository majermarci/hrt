package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	ResponseHeaders http.Header
	RequestHeaders  http.Header
	ResponseBody    string
	RequestName     string
	StatusCode      string
	TLSInfo         *tls.ConnectionState
	Method          string
	URL             string
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

	// Check if both Basic Auth and Bearer Token are specified
	if ep.BasicAuth.Username != "" && ep.BasicAuth.Password != "" && ep.BearerToken != "" {
		fmt.Printf("Error: Both Basic Auth and Bearer Token are specified for endpoint %s.\nPlease select only one method of authentication.\n", test)
		os.Exit(1)
	}

	// Add basic authentication to the request if specified
	if ep.BasicAuth.Username != "" && ep.BasicAuth.Password != "" {
		req.SetBasicAuth(ep.BasicAuth.Username, ep.BasicAuth.Password)
	}

	// Add bearer token to the request if specified
	if ep.BearerToken != "" {
		req.Header.Add("Authorization", "Bearer "+ep.BearerToken)
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

	// Return the response
	return Response{
		ResponseHeaders: resp.Header,
		RequestHeaders:  req.Header,
		ResponseBody:    string(body),
		RequestName:     test,
		StatusCode:      resp.Status,
		TLSInfo:         resp.TLS,
		Method:          ep.Method,
		URL:             ep.URL,
	}
}
