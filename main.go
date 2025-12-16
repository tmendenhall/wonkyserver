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
		file  string
		port  int
		wonky int
		help  bool
	)

	flag.StringVar(&file, "file", "", "Specify a configuration file (required)")
	flag.StringVar(&file, "f", "", "Specify a configuration file (required)")
	flag.IntVar(&port, "port", 8888, "The port that the server will listen to (default: 8888)")
	flag.IntVar(&port, "p", 8888, "The port that the server will listen to (default: 8888)")
	flag.IntVar(&wonky, "wonky", 0, "Percentage (1-100) likelihood of random error/delay/slow behavior (default: 0)")
	flag.IntVar(&wonky, "w", 0, "Percentage (1-100) likelihood of random error/delay/slow behavior (default: 0)")
	flag.BoolVar(&help, "help", false, "Provide a list of all commands and arguments")
	flag.BoolVar(&help, "h", false, "Provide a list of all commands and arguments")

	flag.Parse()

	// Handle help flag
	if help {
		fmt.Println("WonkyServer - A configurable HTTP mock server")
		fmt.Println("\nUsage:")
		fmt.Println("  wonkyserver --file <config.json> [--port <port>] [--wonky <percentage>]")
		fmt.Println("\nArguments:")
		fmt.Println("  --file, -f     Specify a configuration file (required)")
		fmt.Println("  --port, -p     The port that the server will listen to (default: 8888)")
		fmt.Println("  --wonky, -w    Percentage (1-100) likelihood of random error/delay/slow (default: 0)")
		fmt.Println("  --help, -h     Provide a list of all commands and arguments")
		os.Exit(0)
	}

	// Validate flags
	if wonky < 0 || wonky > 100 {
		log.Fatal("Error: --wonky must be between 0 and 100")
	}

	if file == "" {
		log.Fatal("Error: --file flag is required")
	}

	// Load configuration
	config, err := LoadConfig(file)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	log.Printf("Configuration loaded successfully with %d endpoints", len(config.Endpoints))

	if wonky > 0 {
		log.Printf("Wonky mode enabled with %d%% likelihood", wonky)
	}

	// Start server
	if err := StartServer(config, port, wonky); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
