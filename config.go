package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Endpoint represents a single endpoint configuration
type Endpoint struct {
	Verb     string   `json:"verb"`
	URL      string   `json:"url"`
	Code     string   `json:"code"`
	Response string   `json:"response"`
	Headers  []string `json:"headers,omitempty"`
}

// Config represents the overall configuration
type Config struct {
	Endpoints []Endpoint `json:"endpoints"`
}

// LoadConfig loads and validates the configuration from a JSON file
func LoadConfig(filename string) (*Config, error) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file does not exist: %s", filename)
	}

	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	// Check for empty array
	if len(config.Endpoints) == 0 {
		return fmt.Errorf("configuration must contain at least one endpoint")
	}

	// Validate each endpoint
	for i, endpoint := range config.Endpoints {
		if endpoint.Verb == "" {
			return fmt.Errorf("endpoint %d: verb is required", i)
		}
		if endpoint.URL == "" {
			return fmt.Errorf("endpoint %d: url is required", i)
		}
		if endpoint.Code == "" {
			return fmt.Errorf("endpoint %d: code is required", i)
		}
		if endpoint.Response == "" {
			return fmt.Errorf("endpoint %d: response is required", i)
		}
	}

	return nil
}
