package main

import (
	"bufio"
	"embed"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

//go:embed example_config.yaml
var files embed.FS

type BasicAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Endpoint struct {
	URL         string            `yaml:"url"`
	Method      string            `yaml:"method"`
	Body        string            `yaml:"body"`
	Headers     map[string]string `yaml:"headers"`
	BasicAuth   BasicAuth         `yaml:"basic_auth"`
	BearerToken string            `yaml:"bearer_token"`
}

func checkConfig(confFile *string, createGlobal *bool) {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
	}

	globalConfFile := filepath.Join(usr.HomeDir, ".config", cmdName, "config.yaml")

	if *createGlobal {
		if _, err := os.Stat(globalConfFile); err == nil {
			fmt.Printf("Global config file already exists at %s\n", globalConfFile)
			fmt.Println("Please edit the config file to your needs and try again.")
			os.Exit(0)
		} else {
			err := os.MkdirAll(filepath.Dir(globalConfFile), 0700)
			if err != nil {
				log.Fatalf("Failed to create config directory: %v", err)
			}

			createConfig(globalConfFile)

			fmt.Printf("Created global config file at %s\n", globalConfFile)
			fmt.Println("Please edit the config file to your needs and try again.")
			os.Exit(0)
		}
	}

	if *confFile == cmdName+".yaml" {
		if _, err := os.Stat(*confFile); os.IsNotExist(err) {
			if _, err := os.Stat(strings.TrimSuffix(*confFile, "yaml") + "yml"); err == nil {
				*confFile = strings.TrimSuffix(*confFile, "yaml") + "yml"
			} else {
				if _, err := os.Stat(globalConfFile); err == nil {
					*confFile = globalConfFile
				}
			}
		}
	}
}

func createConfig(file string) error {
	exampleConfig, err := files.ReadFile("example_config.yaml")

	err = os.WriteFile(file, exampleConfig, 0600)
	if err != nil {
		return err
	}

	return nil
}

func loadConfig(file string) (map[string]Endpoint, error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Printf("Config file '%v' does not exist.\nDo you want to create it? (Y/n): ", file)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "" || response == "y" {
			err := createConfig(file)
			if err != nil {
				return nil, err
			}
			fmt.Println("\nDefault config file created.\nPlease modify it according to your needs and run the program again.")
			os.Exit(0)
		} else if response == "n" {
			fmt.Println("\nNo config file created. Check '-h' for options.\nExiting.")
			os.Exit(0)
		} else {
			fmt.Println("\nInvalid response. Exiting.")
			os.Exit(1)
		}
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal the YAML data into a map of Endpoint structs
	var c map[string]Endpoint
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
