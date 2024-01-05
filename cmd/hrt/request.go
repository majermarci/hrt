package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type reqResult struct {
	ResponseHeaders http.Header
	RequestHeaders  http.Header
	ResponseBody    string
	RequestName     string
	StatusCode      string
	TLSInfo         *tls.ConnectionState
	Method          string
	URL             string
}

func createRequest(ep endpoint) (*http.Request, error) {
	if ep.Method == "" {
		ep.Method = "GET"
	}

	// Perform an HTTP request for the endpoint
	req, err := http.NewRequest(ep.Method, ep.URL, strings.NewReader(ep.Body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add headers to the request
	for key, value := range ep.Headers {
		req.Header.Add(key, value)
	}

	// Check if both Basic Auth and Bearer Token are specified
	if (ep.BasicAuth.Username != "" || ep.BasicAuth.Password != "") && ep.BearerToken != "" {
		return nil, fmt.Errorf("both basic auth and bearer token are specified\nPlease select only one method of authentication")
	}

	// Add basic authentication to the request if specified
	if ep.BasicAuth.Username != "" && ep.BasicAuth.Password != "" {
		req.SetBasicAuth(ep.BasicAuth.Username, ep.BasicAuth.Password)
	} else if ep.BearerToken != "" {
		// Add bearer token to the request if specified
		req.Header.Add("Authorization", "Bearer "+ep.BearerToken)
	}

	return req, nil
}

func sendRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request - %w", err)
	}
	return resp, nil
}

func handleResponse(resp *http.Response) (string, error) {
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Try to unmarshal the body into a JSON object
	var jsonObj interface{}
	if err := json.Unmarshal(body, &jsonObj); err == nil {
		// If successful, re-marshal it with indentation
		body, _ = json.MarshalIndent(jsonObj, "", "  ")
	}

	return string(body), nil
}

func runTest(test string, ep endpoint, client *http.Client) (reqResult, error) {
	req, err := createRequest(ep)
	if err != nil {
		return reqResult{}, err
	}

	resp, err := sendRequest(client, req)
	if err != nil {
		return reqResult{}, err
	}

	body, err := handleResponse(resp)
	if err != nil {
		return reqResult{}, err
	}

	return reqResult{
		ResponseHeaders: resp.Header,
		RequestHeaders:  req.Header,
		ResponseBody:    body,
		RequestName:     test,
		StatusCode:      resp.Status,
		TLSInfo:         resp.TLS,
		Method:          ep.Method,
		URL:             ep.URL,
	}, nil
}
