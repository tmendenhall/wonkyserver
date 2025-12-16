package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// StartServer starts the HTTP server with the given configuration
func StartServer(config *Config, port int) error {
	handler := &Handler{config: config}

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on port %d", port)

	return http.ListenAndServe(addr, handler)
}

// Handler handles HTTP requests
type Handler struct {
	config *Config
}

// ServeHTTP handles incoming HTTP requests
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for %s", r.Method, r.URL.Path)

	// Find matching endpoint
	endpoint := h.findEndpoint(r.Method, r.URL.Path)
	if endpoint == nil {
		log.Printf("No matching endpoint found, returning 404")
		http.NotFound(w, r)
		return
	}

	// Get query parameters
	query := r.URL.Query()

	// Handle delay parameter
	if delayStr := query.Get("delay"); delayStr != "" {
		if delay, err := parseDelay(delayStr); err == nil {
			log.Printf("Delaying response by %v", delay)
			time.Sleep(delay)
		} else {
			log.Printf("Invalid delay parameter: %v", err)
		}
	}

	// Determine response code
	statusCode := 200

	// Handle error parameter
	if query.Has("error") {
		statusCode = 500
		log.Printf("Error parameter detected, returning 500")
	} else if query.Has("slow") {
		// Per spec: "if the url matches then return a 405 to fast response code"
		statusCode = 405
		log.Printf("Slow parameter detected, returning 405")
	} else {
		// Parse configured status code
		if code, err := strconv.Atoi(endpoint.Code); err == nil {
			statusCode = code
		}
	}

	// Set headers
	for _, header := range endpoint.Headers {
		// Simple header parsing - assumes format like "application/json"
		// Sets Content-Type by default if no colon present
		if strings.Contains(header, ":") {
			parts := strings.SplitN(header, ":", 2)
			w.Header().Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		} else {
			w.Header().Set("Content-Type", header)
		}
	}

	// Write response
	w.WriteHeader(statusCode)
	w.Write([]byte(endpoint.Response))

	log.Printf("Returned %d status code", statusCode)
}

// findEndpoint finds a matching endpoint based on verb and URL
func (h *Handler) findEndpoint(method, path string) *Endpoint {
	for _, endpoint := range h.config.Endpoints {
		if endpoint.Verb == method && endpoint.URL == path {
			return &endpoint
		}
	}
	return nil
}

// parseDelay parses a delay string like "100m", "10s", "1M"
func parseDelay(delayStr string) (time.Duration, error) {
	if len(delayStr) < 2 {
		return 0, fmt.Errorf("invalid delay format")
	}

	// Extract number and unit
	numStr := delayStr[:len(delayStr)-1]
	unit := delayStr[len(delayStr)-1:]

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, fmt.Errorf("invalid delay number: %w", err)
	}

	switch unit {
	case "m":
		return time.Duration(num) * time.Millisecond, nil
	case "s":
		return time.Duration(num) * time.Second, nil
	case "M":
		return time.Duration(num) * time.Minute, nil
	default:
		return 0, fmt.Errorf("invalid delay unit: %s", unit)
	}
}
