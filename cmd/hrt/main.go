package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	// Define the command-line options
	requestName := flag.String("r", "", "Request to run from config file")
	configFile := flag.String("c", "hrt.yaml", "Specify a config file")
	insecure := flag.Bool("k", false, "Disable certificate validation")
	runAll := flag.Bool("a", false, "Run all tests from config file")
	flag.Parse()

	// Load the config file
	conf, err := loadConfig(*configFile)
	if err != nil {
		fmt.Printf("Failed to load config - %v\n", err)
		os.Exit(1)
	}

	// Create an HTTP client
	var client *http.Client
	if *insecure {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	} else {
		client = &http.Client{}
	}

	// If the -a option is enabled, run all tests
	if *runAll {
		if *requestName != "" {
			fmt.Println("Cannot use -a option with a specific request")
			os.Exit(1)
		}

		for requestName, endpoint := range conf {
			runTest(requestName, endpoint, client)
		}
	} else {
		if *requestName == "" {
			fmt.Println("Please specify a request to run using the '-r' flag")
			os.Exit(1)
		}

		// Find the endpoint for the specified request
		endpoint, ok := conf[*requestName]
		if !ok {
			fmt.Printf("No endpoint found for request '%s'\n", *requestName)
			os.Exit(1)
		}

		runTest(*requestName, endpoint, client)
	}

	// If no flags are specified, print out the available flags
	if flag.NFlag() == 0 {
		fmt.Println("Available flags:")
		flag.PrintDefaults()
	}
}
