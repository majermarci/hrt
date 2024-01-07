package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"testing"
)

func TestPrintTLSInfo(t *testing.T) {
	response := reqResult{
		TLSInfo: &tls.ConnectionState{
			Version:     tls.VersionTLS13,
			CipherSuite: tls.TLS_AES_128_GCM_SHA256,
			ServerName:  "fake.domain",
			PeerCertificates: []*x509.Certificate{
				{
					Subject: pkix.Name{
						CommonName: "fake.domain",
					},
					Issuer: pkix.Name{
						CommonName: "ARealCA",
					},
				},
			},
		},
	}

	// Get output of the TLS info func
	output := captureOutput(func() { printTLSInfo(response) })

	// Check that the output is as expected
	expected := "TLS details for endpoint ''\n  TLS version: 1.3\n  Cipher suite: TLS_AES_128_GCM_SHA256\n  Server name: fake.domain\n  Peer certificate: CN=fake.domain\n  Issuer: CN=ARealCA"

	stringContains(t, output, expected)
}

func TestPrintHeaders(t *testing.T) {
	headers := map[string][]string{
		"Content-Type": {"application/json"},
	}

	// Get the output of a "Test headers" listing
	output := captureOutput(func() { printHeaders("Test", headers) })

	// Check that the output is as expected
	expected := "Test Headers:\n  Content-Type: application/json\n"

	equalsTo(t, output, expected)
}

func TestPrintResponseBody(t *testing.T) {
	// Create a test response
	response := reqResult{
		ResponseBody: "All Good!",
	}

	// Get reponse body output
	output := captureOutput(func() { printResponseBody(response) })

	// Check that the output is as expected
	expected := "\nResponse Body:\nAll Good!\n"

	equalsTo(t, output, expected)
}

func TestPrintResponses(t *testing.T) {
	// Create a buffer for capturing output
	var buf bytes.Buffer
	fmt.Fprint(&buf, "")

	// Create a mock response slice
	responses := []reqResult{
		{
			TLSInfo:         nil,
			RequestHeaders:  map[string][]string{"header1": {"value1"}, "header2": {"value2"}},
			ResponseHeaders: map[string][]string{"header3": {"value3"}, "header4": {"value4"}},
			RequestName:     "Request-Test",
			Method:          "GET",
			URL:             "http://fake.domain",
			StatusCode:      "200",
			ResponseBody:    "Response body test",
		},
	}

	// Set verbose and moreVerbose to true
	verbose = new(bool)
	*verbose = true
	moreVerbose = new(bool)
	*moreVerbose = true

	// Get the output of the responses
	output := captureOutput(func() { printResponses(responses) })

	// Check if the output is as expected
	expected := "Request Headers:\n  header1: value1\n  header2: value2\nResponse Headers:\n  header3: value3\n  header4: value4\n\nRequest: 'Request-Test' - GET http://fake.domain\nStatus: 200\n\nResponse Body:\nResponse body test\n"

	equalsTo(t, output, expected)
}
