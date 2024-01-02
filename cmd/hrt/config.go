package main

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

//go:embed example_config.yaml
var Files embed.FS

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

func createDefaultConfig(file string) error {
	// defaultConfig := map[string]Endpoint{
	// 	"example": {
	// 		URL:    "http://localhost:8080",
	// 		Method: "GET",
	// 		BasicAuth: BasicAuth{
	// 			Username: "username",
	// 			Password: "password",
	// 		},
	// 		BearerToken: "your_token_here",
	// 	},
	// }

	// data, err := yaml.Marshal(&defaultConfig)
	// if err != nil {
	// 	return err
	// }

	exampleConfig, err := Files.ReadFile("example_config.yaml")

	err = os.WriteFile(file, exampleConfig, 0644)
	if err != nil {
		return err
	}

	return nil
}

func loadConfig(file string) (map[string]Endpoint, error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Printf("Config file '%v' does not exist.\nDo you want to create a default one? (Y/n): ", file)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "" || response == "y" {
			err := createDefaultConfig(file)
			if err != nil {
				return nil, err
			}
			fmt.Println("Default config file created. Please modify it according to your needs and run the program again.")
			os.Exit(0)
		} else if response == "n" {
			fmt.Println("No config file created. Exiting.")
			os.Exit(0)
		} else {
			fmt.Println("Invalid response. Exiting.")
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
