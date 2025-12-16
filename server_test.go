package main

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_BasicRequest(t *testing.T) {
	config := &Config{
		Endpoints: []Endpoint{
			{
				Verb:     "GET",
				URL:      "/test",
				Code:     "200",
				Response: `{"status":"ok"}`,
				Headers:  []string{"application/json"},
			},
		},
	}

	handler := &Handler{config: config}
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != `{"status":"ok"}` {
		t.Errorf("Expected body '{\"status\":\"ok\"}', got '%s'", w.Body.String())
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestHandler_NotFound(t *testing.T) {
	config := &Config{
		Endpoints: []Endpoint{
			{
				Verb:     "GET",
				URL:      "/test",
				Code:     "200",
				Response: "{}",
			},
		},
	}

	handler := &Handler{config: config}
	req := httptest.NewRequest("GET", "/notfound", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestHandler_ErrorParameter(t *testing.T) {
	config := &Config{
		Endpoints: []Endpoint{
			{
				Verb:     "GET",
				URL:      "/test",
				Code:     "200",
				Response: "{}",
			},
		},
	}

	handler := &Handler{config: config}
	req := httptest.NewRequest("GET", "/test?error", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestHandler_SlowParameter(t *testing.T) {
	config := &Config{
		Endpoints: []Endpoint{
			{
				Verb:     "GET",
				URL:      "/test",
				Code:     "200",
				Response: "{}",
			},
		},
	}

	handler := &Handler{config: config}
	req := httptest.NewRequest("GET", "/test?slow", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != 405 {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHandler_DelayParameter(t *testing.T) {
	config := &Config{
		Endpoints: []Endpoint{
			{
				Verb:     "GET",
				URL:      "/test",
				Code:     "200",
				Response: "{}",
			},
		},
	}

	handler := &Handler{config: config}

	tests := []struct {
		name          string
		delay         string
		minDuration   time.Duration
		expectedValid bool
	}{
		{"100 milliseconds", "100m", 100 * time.Millisecond, true},
		{"1 second", "1s", 1 * time.Second, true},
		{"invalid", "invalid", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test?delay="+tt.delay, nil)
			w := httptest.NewRecorder()

			start := time.Now()
			handler.ServeHTTP(w, req)
			duration := time.Since(start)

			if tt.expectedValid {
				if duration < tt.minDuration {
					t.Errorf("Expected delay of at least %v, got %v", tt.minDuration, duration)
				}
			}

			if w.Code != 200 {
				t.Errorf("Expected status 200, got %d", w.Code)
			}
		})
	}
}

func TestHandler_MultipleEndpoints(t *testing.T) {
	config := &Config{
		Endpoints: []Endpoint{
			{
				Verb:     "GET",
				URL:      "/users",
				Code:     "200",
				Response: `[{"id":1}]`,
			},
			{
				Verb:     "POST",
				URL:      "/users",
				Code:     "201",
				Response: `{"id":2}`,
			},
			{
				Verb:     "GET",
				URL:      "/posts",
				Code:     "200",
				Response: `[]`,
			},
		},
	}

	handler := &Handler{config: config}

	tests := []struct {
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{"GET", "/users", 200, `[{"id":1}]`},
		{"POST", "/users", 201, `{"id":2}`},
		{"GET", "/posts", 200, `[]`},
		{"DELETE", "/users", 404, "404 page not found\n"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestParseDelay(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
		hasError bool
	}{
		{"100m", 100 * time.Millisecond, false},
		{"10s", 10 * time.Second, false},
		{"1M", 1 * time.Minute, false},
		{"invalid", 0, true},
		{"100", 0, true},
		{"x10m", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := parseDelay(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error for input '%s', got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for input '%s', got: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("Expected duration %v, got %v", tt.expected, result)
				}
			}
		})
	}
}
