# WonkyServer

A configurable HTTP mock server written in Go. WonkyServer responds to HTTP requests based on JSON configuration files, allowing you to simulate various API behaviors including errors, delays, and custom responses.

## Features

- JSON-based endpoint configuration
- Support for all HTTP methods (GET, POST, PUT, DELETE, etc.)
- Configurable status codes and response bodies
- Custom response headers
- Query parameter modifiers:
  - `error`: Force 500 error response
  - `slow`: Force 405 response
  - `delay`: Add response delay (milliseconds, seconds, or minutes)
- Multi-architecture Docker support (AMD64 and ARM64)
- Comprehensive logging to stdout

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/wonkyserver.git
cd wonkyserver

# Initialize Go module
go mod init wonkyserver

# Build the binary
go build -o wonkyserver

# Run
./wonkyserver --file example-config.json
```

### Using Docker

Pull from GitHub Container Registry:

```bash
docker pull ghcr.io/yourusername/wonkyserver:latest
```

## Usage

### Command Line Arguments

```bash
wonkyserver --file <config.json> [--port <port>]
```

**Arguments:**
- `--file`, `-f`: Path to configuration file (required)
- `--port`, `-p`: Server port (default: 8888)
- `--help`, `-h`: Show help message

### Configuration File Format

Create a JSON file with an array of endpoint configurations:

```json
{
  "endpoints": [
    {
      "verb": "GET",
      "url": "/api/users",
      "code": "200",
      "response": "[{\"id\":1,\"name\":\"Alice\"}]",
      "headers": ["application/json"]
    }
  ]
}
```

**Fields:**
- `verb`: HTTP method (GET, POST, PUT, DELETE, etc.) - **required**
- `url`: URL path to match - **required**
- `code`: HTTP status code to return - **required**
- `response`: Response body - **required**
- `headers`: Array of response headers - **optional**

### Query Parameters

Modify endpoint behavior using query parameters:

**Error Simulation:**
```bash
curl http://localhost:8888/api/users?error
# Returns 500 instead of configured status
```

**Slow Response:**
```bash
curl http://localhost:8888/api/users?slow
# Returns 405 status code
```

**Delay Response:**
```bash
# Delay by 100 milliseconds
curl http://localhost:8888/api/users?delay=100m

# Delay by 10 seconds
curl http://localhost:8888/api/users?delay=10s

# Delay by 1 minute
curl http://localhost:8888/api/users?delay=1M
```

## Docker Usage

### Running with Docker

```bash
# Using a local config file
docker run -p 8888:8888 \
  -v $(pwd)/example-config.json:/config.json \
  ghcr.io/yourusername/wonkyserver:latest \
  --file /config.json

# Using a different port
docker run -p 9000:9000 \
  -v $(pwd)/example-config.json:/config.json \
  ghcr.io/yourusername/wonkyserver:latest \
  --file /config.json --port 9000
```

### Building Docker Image Locally

```bash
# Build for your platform
docker build -t wonkyserver .

# Build for multiple platforms
docker buildx build --platform linux/amd64,linux/arm64 -t wonkyserver .
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a specific test
go test -run TestHandler_BasicRequest ./...
```

### Project Structure

```
wonkyserver/
├── main.go           # Entry point and CLI parsing
├── config.go         # Configuration loading and validation
├── server.go         # HTTP server and request handling
├── config_test.go    # Configuration tests
├── server_test.go    # Server handler tests
├── example-config.json
├── CURL_EXAMPLES.md  # cURL test examples
├── Dockerfile
└── .github/
    └── workflows/
        └── build.yml # CI/CD pipeline
```

### Adding New Features

This project uses only Go standard library by design. If you need to add external dependencies, please open an issue first to discuss.

## Examples

See [CURL_EXAMPLES.md](CURL_EXAMPLES.md) for comprehensive testing examples.

### Quick Example

1. Create a config file `my-config.json`:

```json
{
  "endpoints": [
    {
      "verb": "GET",
      "url": "/hello",
      "code": "200",
      "response": "{\"message\":\"Hello, World!\"}",
      "headers": ["application/json"]
    }
  ]
}
```

2. Run the server:

```bash
go run . --file my-config.json
```

3. Test it:

```bash
curl http://localhost:8888/hello
# Output: {"message":"Hello, World!"}

curl http://localhost:8888/hello?error
# Output: {"message":"Hello, World!"} with HTTP 500 status

curl http://localhost:8888/hello?delay=2s
# Output: {"message":"Hello, World!"} after 2 second delay
```

## GitHub Actions

The project includes a GitHub Actions workflow that:
- Runs tests on every push and pull request
- Builds Docker images for AMD64 and ARM64
- Publishes images to GitHub Container Registry
- Tags images based on branch, tag, and commit SHA

To enable, ensure GitHub Actions has package write permissions in your repository settings.

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
