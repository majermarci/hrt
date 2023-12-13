package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
)

var (
	confFile     = flag.String("c", "hrt.yaml", "Specify a config file")
	requestName  = flag.String("r", "", "Request to run from config file")
	runAll       = flag.Bool("a", false, "Run all tests from config file")
	insecure     = flag.Bool("k", false, "Disable certificate validation")
	tableOutput  = flag.Bool("t", false, "Enable table output")
	allResponses []Response
)

func main() {
	flag.Parse()

	// Load the configuration file
	conf, err := loadConfig(*confFile)
	if err != nil {
		fmt.Printf("Failed to load configuration file: %v\n", err)
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
			response := runTest(requestName, endpoint, client)
			allResponses = append(allResponses, response)
		}
		if *tableOutput {
			printTable(allResponses)
		} else {
			printResponses(allResponses)
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

		response := runTest(*requestName, endpoint, client)
		allResponses = append(allResponses, response)
		if *tableOutput {
			printTable(allResponses)
		} else {
			printResponses(allResponses)
		}
	}

	// If no flags are specified, print out the available flags
	if flag.NFlag() == 0 {
		fmt.Println("Available flags:")
		flag.PrintDefaults()
	}
}
