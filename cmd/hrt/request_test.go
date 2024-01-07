package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateRequest(t *testing.T) {
	ep := endpoint{
		Method: "GET",
		URL:    "http://fake.domain",
		Headers: map[string]string{
			"X-Header-1": "Value1",
		},
		BasicAuth: basicAuth{
			Username: "test",
			Password: "test",
		},
	}

	// Make the request
	req, err := createRequest(ep)

	// Check that there is no error
	nilError(t, err)

	// Check that the request method is as expected
	equalsTo(t, req.Method, http.MethodGet)

	// Check that the request URL is as expected
	equalsTo(t, req.URL.String(), ep.URL)
}

func TestSendRequest(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("OK"))
	}))
	// Close the server when test finishes
	defer server.Close()

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)

	// Create a test client
	client := &http.Client{}

	// Call the function with the test client and request
	resp, err := sendRequest(client, req)

	// Check that there is no error
	nilError(t, err)

	// Check that the response status code is as expected
	equalsTo(t, resp.StatusCode, http.StatusOK)
}

func TestHandleResponse(t *testing.T) {
	// Create a mock response with a JSON body
	resp := &http.Response{
		Body:       io.NopCloser(strings.NewReader(`{"message": "OK"}`)),
		StatusCode: http.StatusOK,
	}

	// Call the function with the test response
	body, err := handleResponse(resp)

	// Check that there is no error
	nilError(t, err)

	// Check that the response body is as expected
	expectedBody := `{
  "message": "OK"
}`

	equalsTo(t, body, expectedBody)
}

func TestRunTest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "ok"}`))
	}))
	defer ts.Close()

	// Create mock client and endpoint
	client := ts.Client()
	ep := endpoint{
		Method: "GET",
		URL:    ts.URL,
	}

	// Run the test
	result, err := runTest("Test1", ep, client)

	expectedBody := `{
  "message": "ok"
}`

	// Check that there is no error
	nilError(t, err)

	// Check the results
	equalsTo(t, result.RequestName, "Test1")
	equalsTo(t, result.TLSInfo, nil)
	equalsTo(t, result.Method, http.MethodGet)
	equalsTo(t, result.URL, ts.URL)
	equalsTo(t, result.StatusCode, "200 OK")
	equalsTo(t, result.ResponseBody, expectedBody)
}
