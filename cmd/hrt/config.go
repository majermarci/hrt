package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

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

func loadConfig(file string) (map[string]Endpoint, error) {
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
