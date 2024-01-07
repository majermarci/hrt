package main

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	file, err := os.CreateTemp("", "config")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Call the function with the temp file
	err = createConfig(file.Name())
	if err != nil {
		t.Errorf("error creating config: %v", err)
	}

	// Check that the file now exists
	if !checkConfigExists(file.Name()) {
		t.Errorf("file does not exist")
	}

	// Check that the function returns false for a non-existent file
	if checkConfigExists("non_existent_file") {
		t.Errorf("expected false but got true")
	}
}

func TestCreateGlobalConfigFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "config")
	if err != nil {
		t.Fatalf("failed to create tempdir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Call the function with a file in the temporary directory
	err = createGlobalConfigFile(dir + "/config.yaml")
	if err != nil {
		t.Errorf("error creating global config: %v", err)
	}

	// Check that the file now exists
	if !checkConfigExists(dir + "/config.yaml") {
		t.Errorf("file does not exist")
	}
}

func TestLoadConfig(t *testing.T) {
	file, err := os.CreateTemp("", "config")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	// Write an example config in the temp file
	_, err = file.WriteString(`test:
  url: http://example.com
  method: GET
`)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	// Call the function with the temp file
	config, err := loadConfig(file.Name())
	if err != nil {
		t.Errorf("error loading config: %v", err)
	}

	// Check that the config is as expected
	if config["test"].URL != "http://example.com" || config["test"].Method != "GET" {
		t.Errorf("expected config with URL 'http://example.com' and method 'GET', got config with URL '%s' and method '%s'", config["test"].URL, config["test"].Method)
	}
}

func TestListAllRequests(t *testing.T) {
	file, err := os.CreateTemp("", "config")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Write an example config in the temp file
	_, err = file.WriteString("test1:\n  url: http://localhost:8080\ntest2:\n  url: http://localhost:8080")
	if err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	// Get the output of the listing...
	output := captureOutput(func() {
		err = listAllRequests(file.Name())
		if err != nil {
			t.Errorf("error listing requests %v", err)
		}
	})

	equalsTo(t, output, "test1\ntest2\n")
}
