package main

import (
	"crypto/tls"
	"fmt"
	"strings"
)

// Helper function to print responses
func printResponses(responses []Response) {
	if *verbose {
		fmt.Printf("Used config file: %s\n\n", *confFile)
	}

	for _, response := range responses {
		// Print TLS details
		if *verbose && response.TLSInfo != nil {
			fmt.Printf("TLS details for endpoint '%s'\n", response.RequestName)
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

		// Print all request headers if verbose is enabled
		if *verbose && len(response.RequestHeaders) > 0 {
			fmt.Println("Request Headers:")
			for key, values := range response.RequestHeaders {
				for _, value := range values {
					fmt.Printf("  %s: %s\n", key, value)
				}
			}
		}

		// Print response headers if verbose is enabled
		if *verbose && len(response.ResponseHeaders) > 0 {
			fmt.Println("Response Headers:")
			for key, values := range response.ResponseHeaders {
				for _, value := range values {
					fmt.Printf("  %s: %s\n", key, value)
				}
			}
		}

		fmt.Printf("\nRequest: '%v' - %v %v\n", response.RequestName, response.Method, response.URL)
		fmt.Printf("Status: %v\n", response.StatusCode)
		if response.ResponseBody != "" {
			fmt.Printf("\nResponse Body:\n %v\n", response.ResponseBody)
		}
	}
}

func printTable(responses []Response) {
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

	// Print the table header
	printLine(testWidth, statusWidth, bodyWidth, "╭", "┬", "╮")
	fmt.Printf("│ %-*s │ %-*s │ %-*s │\n", testWidth, "Request", statusWidth, "Status Code", bodyWidth, "Response Body")
	printLine(testWidth, statusWidth, bodyWidth, "├", "┼", "┤")

	// For each response, print a table row
	for i, response := range responses {
		lines := strings.Split(response.ResponseBody, "\n")
		for j, line := range lines {
			if j == 0 {
				fmt.Printf("│ %-*s │ %-*s │ %-*s │\n", testWidth, response.RequestName, statusWidth, response.StatusCode, bodyWidth, line)
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
