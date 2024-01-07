package main

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

//go:embed example_config.yaml
var files embed.FS

type basicAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type endpoint struct {
	URL         string            `yaml:"url"`
	Method      string            `yaml:"method"`
	Body        string            `yaml:"body"`
	Headers     map[string]string `yaml:"headers"`
	BasicAuth   basicAuth         `yaml:"basic_auth"`
	BearerToken string            `yaml:"bearer_token"`
}

const exampleConfigFile = "example_config.yaml"

func createConfig(file string) error {
	exampleConfig, err := files.ReadFile(exampleConfigFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(file, exampleConfig, 0600)
	if err != nil {
		return err
	}

	return nil
}

func checkConfigExists(confFile string) bool {
	_, err := os.Stat(confFile)
	return !os.IsNotExist(err)
}

func createGlobalConfigFile(globalConfFile string) error {
	if checkConfigExists(globalConfFile) {
		return fmt.Errorf("file already exists at %s", globalConfFile)
	}

	err := os.MkdirAll(filepath.Dir(globalConfFile), 0700)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	err = createConfig(globalConfFile)
	if err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}

	return nil
}

func getUserInput(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.ToLower(strings.TrimSpace(response)), nil
}

func checkConfig(confFile *string, createGlobal *bool) {
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Failed to get current user: %v\n", err)
		os.Exit(1)
	}

	globalConfFile := filepath.Join(usr.HomeDir, ".config", cmdName, "config.yaml")

	if *createGlobal {
		err := createGlobalConfigFile(globalConfFile)
		if err != nil {
			fmt.Printf("Failed to create global config file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Created global config file at %s\n", globalConfFile)
		fmt.Println("Please edit the config file to your needs and try again.")
		os.Exit(0)
	}

	if *confFile == cmdName+".yaml" {
		if !checkConfigExists(*confFile) {
			if checkConfigExists(strings.TrimSuffix(*confFile, "yaml") + "yml") {
				*confFile = strings.TrimSuffix(*confFile, "yaml") + "yml"
			} else {
				if checkConfigExists(globalConfFile) {
					*confFile = globalConfFile
				}
			}
		}
	}
}

func loadConfig(file string) (map[string]endpoint, error) {
	if !checkConfigExists(file) {
		response, err := getUserInput(fmt.Sprintf("Config file '%v' does not exist.\nDo you want to create it? (Y/n): ", file))
		if err != nil {
			return nil, err
		}

		if response == "" || response == "y" {
			err := createConfig(file)
			if err != nil {
				return nil, err
			}
			fmt.Println("\nDefault config file created.\nPlease modify it according to your needs and run the program again.")
			os.Exit(0)
		} else if response == "n" {
			fmt.Println("\nNo config file created. Check '-h' for options")
			os.Exit(0)
		} else {
			return nil, fmt.Errorf("invalid response")
		}
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal the YAML data into a map of Endpoint structs
	var c map[string]endpoint
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func listAllRequests(file string) error {
	if !checkConfigExists(file) {
		return fmt.Errorf("config file '%v' does not exist", file)
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	// Unmarshal the YAML data into a map of Endpoint structs
	var c map[string]endpoint
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return err
	}

	// Print all available requests
	for request := range c {
		fmt.Println(request)
	}

	return nil
}
