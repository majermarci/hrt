package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

type Endpoint struct {
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Body    string            `yaml:"body"`
	Headers map[string]string `yaml:"headers"`
}

func main() {
	// Define the command-line options
	configFile := flag.String("c", "hrt.yaml", "Specify a config file")
	runAll := flag.Bool("a", false, "Run all tests from config file")
	requestName := flag.String("r", "", "Request to run from config file")
	insecure := flag.Bool("k", false, "Disable certificate validation")
	flag.Parse()

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

	// Read the YAML file
	data, err := os.ReadFile(*configFile)
	if err != nil {
		fmt.Printf("Failed to load default config - %v\n", err)
		os.Exit(1)
	}

	// Unmarshal the YAML data into a map of Endpoint structs
	var config map[string]Endpoint
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Failed to unmarshal config - %v\n", err)
		os.Exit(1)
	}

	// If the -a option is enabled, run all tests
	if *runAll {
		if *requestName != "" {
			fmt.Println("Cannot use -a option with a specific request")
			os.Exit(1)
		}

		for requestName, endpoint := range config {
			runTest(requestName, endpoint, client)
		}
	} else {
		if *requestName == "" {
			fmt.Println("Please specify a request to run using the '-r' flag")
			os.Exit(1)
		}

		// Find the endpoint for the specified request
		endpoint, ok := config[*requestName]
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
