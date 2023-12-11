package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Endpoint struct {
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Body    string            `yaml:"body"`
	Headers map[string]string `yaml:"headers"`
}

func loadConfig(file string) (map[string]Endpoint, error) {
	data, err := os.ReadFile(*&file)
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
