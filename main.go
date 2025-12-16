package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// Define command line flags
	var (
		file string
		port int
		help bool
	)

	flag.StringVar(&file, "file", "", "Specify a configuration file (required)")
	flag.StringVar(&file, "f", "", "Specify a configuration file (required)")
	flag.IntVar(&port, "port", 8888, "The port that the server will listen to (default: 8888)")
	flag.IntVar(&port, "p", 8888, "The port that the server will listen to (default: 8888)")
	flag.BoolVar(&help, "help", false, "Provide a list of all commands and arguments")
	flag.BoolVar(&help, "h", false, "Provide a list of all commands and arguments")

	flag.Parse()

	// Handle help flag
	if help {
		fmt.Println("WonkyServer - A configurable HTTP mock server")
		fmt.Println("\nUsage:")
		fmt.Println("  wonkyserver --file <config.json> [--port <port>]")
		fmt.Println("\nArguments:")
		fmt.Println("  --file, -f    Specify a configuration file (required)")
		fmt.Println("  --port, -p    The port that the server will listen to (default: 8888)")
		fmt.Println("  --help, -h    Provide a list of all commands and arguments")
		os.Exit(0)
	}

	// Validate required flags
	if file == "" {
		log.Fatal("Error: --file flag is required")
	}

	// Load configuration
	config, err := LoadConfig(file)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	log.Printf("Configuration loaded successfully with %d endpoints", len(config.Endpoints))

	// Start server
	if err := StartServer(config, port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
