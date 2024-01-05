package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	cmdName      = filepath.Base(os.Args[0])
	confFile     = flag.String("c", cmdName+".yaml", "Specify a config file")
	certFile     = flag.String("cert", "", "Path to the TLS certificate file")
	keyFile      = flag.String("key", "", "Path to the TLS private key file")
	caCertFile   = flag.String("cacert", "", "Path to the CA certificate file")
	requestName  = flag.String("r", "", "Request to run from config file")
	timeout      = flag.Int("t", 10, "Timeout for the HTTP client in seconds")
	listRequests = flag.Bool("l", false, "List all available requests in the current config file")
	createGlobal = flag.Bool("g", false, "Create a global config file")
	runAll       = flag.Bool("a", false, "Run all tests from config file")
	insecure     = flag.Bool("k", false, "Disable certificate validation")
	tableOutput  = flag.Bool("table", false, "Enable table output (experimental)")
	verbose      = flag.Bool("v", false, "Enable verbose Request details")
	moreVerbose  = flag.Bool("vv", false, "Enable verbose Request and TLS details")
	version      = flag.Bool("version", false, "Print the version")
	allResponses []reqResult
	commitID     string
)

const (
	appVersion = "v0.4.2"
)

func main() {
	flag.Parse()

	// If the -version flag is provided and no other flags are provided, print the version and exit
	if *version && flag.NFlag() == 1 {
		if commitID != "" {
			fmt.Printf("%s, commit %s\n", appVersion, commitID)
		} else {
			fmt.Println(appVersion)
		}
		os.Exit(0)
	}

	// // If no flags are specified, print out the available flags
	// if flag.NFlag() == 0 {
	// 	fmt.Println("Available flags:")
	// 	flag.PrintDefaults()
	// }

	checkConfig(confFile, createGlobal)

	// Load the configuration file
	conf, err := loadConfig(*confFile)
	if err != nil {
		fmt.Printf("Failed to load configuration file: %v\n", err)
		os.Exit(1)
	}

	if *verbose || *moreVerbose {
		fmt.Printf("Used config file: %s\n", *confFile)
	}

	if *listRequests {
		listAllRequests(*confFile)
		os.Exit(0)
	}

	// Load system CA pool
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		fmt.Printf("Failed to load system CA pool: %v", err)
		os.Exit(1)
	}

	// Append additional CA certs from file if specified
	// var caCertPool *x509.CertPool
	if *caCertFile != "" {
		caCert, err := os.ReadFile(*caCertFile)
		if err != nil {
			fmt.Printf("Failed to read CA certificate: %v", err)
			os.Exit(1)
		}

		// caCertPool = x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			fmt.Println("Failed to append CA certificate")
			os.Exit(1)
		}
	}

	// Create a base transport with common settings
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: *insecure,
			RootCAs:            caCertPool,
		},
	}

	// Load certificate if certFile and keyFile are provided
	if *certFile != "" && *keyFile != "" {
		cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
		if err != nil {
			fmt.Printf("Failed to load key pair: %v", err)
			os.Exit(1)
		}
		transport.TLSClientConfig.Certificates = []tls.Certificate{cert}
	}

	// Create an HTTP client with the configured transport
	client := &http.Client{
		Timeout:   time.Duration(*timeout) * time.Second,
		Transport: transport,
	}

	// Check if a request name is provided when -a option is enabled
	if *runAll && *requestName != "" {
		fmt.Println("Cannot use -a option with a specific request")
		os.Exit(1)
	}

	// Check if a request name is not provided when -a option is disabled
	if !*runAll && *requestName == "" {
		fmt.Println("Please specify a request to run using the '-r' flag")
		fmt.Println("Check the '-h' flag for additional help.")
		os.Exit(1)
	}

	// If the -a option is enabled, run all tests, else run the specified test
	requests := make(map[string]endpoint)
	if *runAll {
		requests = conf
	} else {
		endpoint, ok := conf[*requestName]
		if !ok {
			fmt.Printf("No endpoint found for request '%s'\n", *requestName)
			os.Exit(1)
		}
		requests[*requestName] = endpoint
	}

	// Run tests for all requests and collect responses
	for requestName, endpoint := range requests {
		response, err := runTest(requestName, endpoint, client)
		if err != nil {
			fmt.Printf("Failed to run test '%s': %v\n", requestName, err)
			os.Exit(1)
		}
		allResponses = append(allResponses, response)
	}

	// Print responses
	if *tableOutput {
		printTable(allResponses)
	} else {
		printResponses(allResponses)
	}
}
