package main

import (
	"os"
	"testing"
)

func TestLoadConfig_ValidConfig(t *testing.T) {
	// Create a temporary config file
	content := `{
		"endpoints": [
			{
				"verb": "GET",
				"url": "/test",
				"code": "200",
				"response": "{\"status\":\"ok\"}",
				"headers": ["application/json"]
			}
		]
	}`

	tmpfile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Load config
	config, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Validate config
	if len(config.Endpoints) != 1 {
		t.Fatalf("Expected 1 endpoint, got %d", len(config.Endpoints))
	}

	endpoint := config.Endpoints[0]
	if endpoint.Verb != "GET" {
		t.Errorf("Expected verb 'GET', got '%s'", endpoint.Verb)
	}
	if endpoint.URL != "/test" {
		t.Errorf("Expected url '/test', got '%s'", endpoint.URL)
	}
	if endpoint.Code != "200" {
		t.Errorf("Expected code '200', got '%s'", endpoint.Code)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/file.json")
	if err == nil {
		t.Fatal("Expected error for missing file, got nil")
	}
}

func TestLoadConfig_EmptyArray(t *testing.T) {
	content := `{"endpoints": []}`

	tmpfile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = LoadConfig(tmpfile.Name())
	if err == nil {
		t.Fatal("Expected error for empty array, got nil")
	}
}

func TestLoadConfig_MissingRequiredFields(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "missing verb",
			content: `{
				"endpoints": [
					{"url": "/test", "code": "200", "response": "{}"}
				]
			}`,
		},
		{
			name: "missing url",
			content: `{
				"endpoints": [
					{"verb": "GET", "code": "200", "response": "{}"}
				]
			}`,
		},
		{
			name: "missing code",
			content: `{
				"endpoints": [
					{"verb": "GET", "url": "/test", "response": "{}"}
				]
			}`,
		},
		{
			name: "missing response",
			content: `{
				"endpoints": [
					{"verb": "GET", "url": "/test", "code": "200"}
				]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "config*.json")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.content)); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			_, err = LoadConfig(tmpfile.Name())
			if err == nil {
				t.Fatalf("Expected error for %s, got nil", tt.name)
			}
		})
	}
}

func TestLoadConfig_OptionalHeaders(t *testing.T) {
	content := `{
		"endpoints": [
			{
				"verb": "GET",
				"url": "/test",
				"code": "200",
				"response": "{}"
			}
		]
	}`

	tmpfile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	config, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error when headers are missing, got: %v", err)
	}

	if len(config.Endpoints[0].Headers) != 0 {
		t.Errorf("Expected 0 headers, got %d", len(config.Endpoints[0].Headers))
	}
}
