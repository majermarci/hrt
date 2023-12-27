package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"net/http"
	"os"
)

var (
	confFile     = flag.String("c", "hrt.yaml", "Specify a config file")
	certFile     = flag.String("cert", "", "Path to the TLS certificate file")
	keyFile      = flag.String("key", "", "Path to the TLS private key file")
	caCertFile   = flag.String("cacert", "", "Path to the CA certificate file")
	requestName  = flag.String("r", "", "Request to run from config file")
	runAll       = flag.Bool("a", false, "Run all tests from config file")
	insecure     = flag.Bool("k", false, "Disable certificate validation")
	tableOutput  = flag.Bool("t", false, "Enable table output")
	version      = flag.Bool("v", false, "Print the version")
	allResponses []Response
)

const (
	appVersion = "0.2.3"
)

func main() {
	flag.Parse()

	// If the -v flag is provided and no other flags are provided, print the version and exit
	if *version && flag.NFlag() == 1 {
		fmt.Println(appVersion)
		os.Exit(0)
	}

	// If no flags are specified, print out the available flags
	if flag.NFlag() == 0 {
		fmt.Println("Available flags:")
		flag.PrintDefaults()
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

	// Create an HTTP client
	var client *http.Client
	if *insecure {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	} else if *certFile != "" && *keyFile != "" {
		cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
		if err != nil {
			fmt.Printf("Failed to load key pair: %v", err)
			os.Exit(1)
		}

		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates: []tls.Certificate{cert},
					RootCAs:      caCertPool,
				},
			},
		}
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: caCertPool,
				},
			},
		}
	}

	// Load the configuration file
	conf, err := loadConfig(*confFile)
	if err != nil {
		fmt.Printf("Failed to load configuration file: %v\n", err)
		os.Exit(1)
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
}
