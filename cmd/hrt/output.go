package main

import (
	"crypto/tls"
	"fmt"
	"strings"

	"golang.org/x/term"
)

func printTLSInfo(response reqResult) {
	fmt.Printf("\nTLS details for endpoint '%s'\n", response.RequestName)
	tlsVersion := map[uint16]string{
		tls.VersionTLS10: "1.0",
		tls.VersionTLS11: "1.1",
		tls.VersionTLS12: "1.2",
		tls.VersionTLS13: "1.3",
	}
	fmt.Printf("  TLS version: %s\n", tlsVersion[response.TLSInfo.Version])

	tlsCipherSuite := map[uint16]string{
		tls.TLS_AES_128_GCM_SHA256:                        "TLS_AES_128_GCM_SHA256",
		tls.TLS_AES_256_GCM_SHA384:                        "TLS_AES_256_GCM_SHA384",
		tls.TLS_CHACHA20_POLY1305_SHA256:                  "TLS_CHACHA20_POLY1305_SHA256",
		tls.TLS_RSA_WITH_RC4_128_SHA:                      "TLS_RSA_WITH_RC4_128_SHA",
		tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA:                 "TLS_RSA_WITH_3DES_EDE_CBC_SHA",
		tls.TLS_RSA_WITH_AES_128_CBC_SHA:                  "TLS_RSA_WITH_AES_128_CBC_SHA",
		tls.TLS_RSA_WITH_AES_256_CBC_SHA:                  "TLS_RSA_WITH_AES_256_CBC_SHA",
		tls.TLS_RSA_WITH_AES_128_CBC_SHA256:               "TLS_RSA_WITH_AES_128_CBC_SHA256",
		tls.TLS_RSA_WITH_AES_128_GCM_SHA256:               "TLS_RSA_WITH_AES_128_GCM_SHA256",
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384:               "TLS_RSA_WITH_AES_256_GCM_SHA384",
		tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA:              "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA",
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA:          "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA:          "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA",
		tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA:                "TLS_ECDHE_RSA_WITH_RC4_128_SHA",
		tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA:           "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA",
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA:            "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA:            "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256:       "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256",
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256:         "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256:         "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256:       "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384:         "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384:       "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256:   "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256",
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256: "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256",
	}

	fmt.Printf("  Cipher suite: %s\n", tlsCipherSuite[response.TLSInfo.CipherSuite])
	fmt.Printf("  Server name: %s\n", response.TLSInfo.ServerName)
	if response.TLSInfo.PeerCertificates != nil && len(response.TLSInfo.PeerCertificates) > 0 {
		fmt.Printf("  Peer certificate: %v\n", response.TLSInfo.PeerCertificates[0].Subject)
		fmt.Printf("  Issuer: %v\n\n", response.TLSInfo.PeerCertificates[0].Issuer)
	}
}

func printHeaders(headerType string, headers map[string][]string) {
	fmt.Println(headerType + " Headers:")
	for key, values := range headers {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}
}

func printResponseBody(response reqResult) {
	if response.ResponseBody != "" {
		fmt.Println(response.ResponseBody)
	}
}

func printResponses(responses []reqResult) {
	for _, response := range responses {
		if *moreVerbose && response.TLSInfo != nil {
			printTLSInfo(response)
		}

		if *moreVerbose && len(response.RequestHeaders) > 0 {
			printHeaders("Request", response.RequestHeaders)
		}

		// If the method is HEAD or moreVerbose is true, print the response headers
		if (response.Method == "HEAD" || response.Method == "OPTIONS") || (*moreVerbose && len(response.ResponseHeaders) > 0) {
			printHeaders("Response", response.ResponseHeaders)
		}

		if *verbose || *moreVerbose {
			fmt.Printf("\nRequest: '%v' - %v %v\n", response.RequestName, response.Method, response.URL)
			fmt.Printf("Status: %v\n", response.StatusCode)
			fmt.Println("\nResponse Body:")
		}

		printResponseBody(response)
	}
}

func printTable(responses []reqResult) {
	// Get the current terminal size
	termWidth, _, err := term.GetSize(0)
	if err != nil {
		// If there's an error, default to a reasonable terminal width
		termWidth = 80
	}

	// Define the maximum width for each column
	maxTestWidth := 25
	maxStatusWidth := 35
	// maxBodyWidth := 60

	// Determine the width of each column
	testWidth := len("Request")
	statusWidth := len("Status Code")
	bodyWidth := len("Response Body")

	// Find the longest value for each column
	for _, response := range responses {
		if len(response.RequestName) > testWidth {
			testWidth = len(response.RequestName)
		}
		if len(response.StatusCode) > statusWidth {
			statusWidth = len(response.StatusCode)
		}
		lines := strings.Split(response.ResponseBody, "\n")
		for _, line := range lines {
			if len(line) > bodyWidth {
				bodyWidth = len(line)
			}
		}
	}

	// Calculate the maximum width for the body column
	maxBodyWidth := termWidth - testWidth - statusWidth - 10

	// Limit the width of each column to the maximum width
	testWidth = min(testWidth, maxTestWidth)
	statusWidth = min(statusWidth, maxStatusWidth)
	bodyWidth = min(bodyWidth, maxBodyWidth)

	// Print the table header
	printLine(testWidth, statusWidth, bodyWidth, "╭", "┬", "╮")
	fmt.Printf("│ %-*s │ %-*s │ %-*s │\n", testWidth, "Request", statusWidth, "Status Code", bodyWidth, "Response Body")
	printLine(testWidth, statusWidth, bodyWidth, "├", "┼", "┤")

	// For each response, print a table row
	for i, response := range responses {
		lines := strings.Split(response.ResponseBody, "\n")
		for j, line := range lines {
			// Truncate the line if it's too long
			line = truncate(line, maxBodyWidth)
			if j == 0 {
				fmt.Printf("│ %-*s │ %-*s │ %-*s │\n", testWidth, truncate(response.RequestName, maxTestWidth), statusWidth, truncate(response.StatusCode, maxStatusWidth), bodyWidth, line)
			} else {
				fmt.Printf("│ %-*s │ %-*s │ %-*s │\n", testWidth, "", statusWidth, "", bodyWidth, line)
			}
		}
		// Only print the line separator if it's not the last response
		if i < len(responses)-1 {
			printLine(testWidth, statusWidth, bodyWidth, "├", "┼", "┤")
		}
	}

	// Print the table footer
	printLine(testWidth, statusWidth, bodyWidth, "╰", "┴", "╯")
}

func printLine(testWidth, statusWidth, bodyWidth int, start, middle, end string) {
	fmt.Print(start)
	fmt.Print(strings.Repeat("─", testWidth+2))
	fmt.Print(middle)
	fmt.Print(strings.Repeat("─", statusWidth+2))
	fmt.Print(middle)
	fmt.Print(strings.Repeat("─", bodyWidth+2))
	fmt.Println(end)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func truncate(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}
