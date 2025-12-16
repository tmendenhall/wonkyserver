# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

WonkyServer is a configurable HTTP mock server written in Go. It responds to HTTP requests based on JSON configuration files, allowing users to simulate various API behaviors including errors, delays, and custom responses.

## Architecture

### Core Components

- **Configuration Loading**: JSON-based configuration parser that defines endpoint behaviors (verb, URL, response code, body, headers)
- **HTTP Server**: Standard library HTTP server that matches incoming requests against configured endpoints
- **Request Matching**: Routes requests by matching HTTP verb and URL path, returns configured responses
- **Query Parameter Handlers**: Modifies response behavior based on special query parameters:
  - `error`: Returns 500 instead of configured status code
  - `slow`: Returns 405 (fast response code intentionally incorrect per spec)
  - `delay`: Delays response by specified duration (e.g., `delay=100m`, `delay=10s`, `delay=1M`)
- **Wonky Mode**: Global chaos testing feature that randomly applies error/delay/slow behaviors based on configurable percentage (0-100). Explicit query parameters always override wonky behavior.

### Configuration Format

The application expects a JSON file with an array of endpoint configurations:
```json
{
    [
        {"verb": "GET", "url": "/foo", "code":"200", "response":"{}","headers":["application/json"]}
    ]
}
```

All fields are required except `headers`. Empty arrays or missing files cause startup errors.

### Command Line Interface

Arguments use both long (`--arg`) and short (`-a`) forms, with short forms using first letter or first n letters if conflicts exist:

- `--help` / `-h`: Display all commands and arguments
- `--file` / `-f`: Required. Specify configuration file path
- `--port` / `-p`: Optional. Server port (default: 8888)
- `--wonky` / `-w`: Optional. Percentage (1-100) likelihood of random error/delay/slow behavior (default: 0)

## Development Commands

### Standard Go Commands

```bash
# Run the application
go run . --file config.json --port 8888

# Run with wonky mode (30% chaos)
go run . --file config.json --wonky 30

# Build the binary
go build -o wonkyserver

# Run tests
go test ./...

# Run a single test
go test -run TestName ./...

# Run tests with verbose output
go test -v ./...
```

### Container Operations

The application is designed to run in containers with multi-architecture support (ARM64 and x86).

```bash
# Build container
docker build -t wonkyserver .

# Run container
docker run -p 8888:8888 -v $(pwd)/config.json:/config.json wonkyserver --file /config.json
```

## Implementation Constraints

- **Standard Library Only**: Use only Go standard library imports unless explicitly approved. Ask before adding external dependencies.
- **Logging**: All messages must log to standard output (stdout)
- **Container First**: Application designed for containerized deployment with GitHub Actions CI/CD
- **Multi-arch Support**: Container builds must support both ARM64 and x86_64 architectures
