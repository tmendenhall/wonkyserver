package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// StartServer starts the HTTP server with the given configuration
func StartServer(config *Config, port int, wonkyPercentage int) error {
	handler := &Handler{
		config:          config,
		wonkyPercentage: wonkyPercentage,
	}

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on port %d", port)

	return http.ListenAndServe(addr, handler)
}

// Handler handles HTTP requests
type Handler struct {
	config          *Config
	wonkyPercentage int
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

	// Apply wonky behavior randomly
	wonkyBehavior := applyWonkyBehavior(h.wonkyPercentage)
	if wonkyBehavior != "" {
		log.Printf("Wonky behavior applied: %s", wonkyBehavior)
	}

	// Handle delay parameter (explicit param takes precedence over wonky)
	if delayStr := query.Get("delay"); delayStr != "" {
		if delay, err := parseDelay(delayStr); err == nil {
			log.Printf("Delaying response by %v", delay)
			time.Sleep(delay)
		} else {
			log.Printf("Invalid delay parameter: %v", err)
		}
	} else if wonkyBehavior == "delay5s" {
		log.Printf("Wonky delay: 5 seconds")
		time.Sleep(5 * time.Second)
	}

	// Determine response code
	statusCode := 200

	// Handle error parameter (explicit params take precedence over wonky)
	if query.Has("error") {
		statusCode = 500
		log.Printf("Error parameter detected, returning 500")
	} else if query.Has("slow") {
		// Return 429 Too Many Requests for rate limiting
		statusCode = 429
		log.Printf("Slow parameter detected, returning 429")
	} else if wonkyBehavior == "error" {
		statusCode = 500
		log.Printf("Wonky error: returning 500")
	} else if wonkyBehavior == "slow" {
		statusCode = 429
		log.Printf("Wonky slow: returning 429")
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

// applyWonkyBehavior determines if wonky behavior should be applied
// Returns one of: "error", "slow", "delay5s", or "" (no wonky behavior)
func applyWonkyBehavior(percentage int) string {
	if percentage == 0 {
		return ""
	}

	// Generate random number 1-100
	rand.Seed(time.Now().UnixNano())
	roll := rand.Intn(100) + 1

	if roll <= percentage {
		// Randomly select one of three behaviors
		behaviors := []string{"error", "slow", "delay5s"}
		return behaviors[rand.Intn(3)]
	}

	return ""
}
