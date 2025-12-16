# cURL Test Examples

This document provides example cURL commands to test the WonkyServer endpoints.

## Prerequisites

Start the server with the example configuration:

```bash
go run . --file example-config.json --port 8888
```

Or using the compiled binary:

```bash
./wonkyserver --file example-config.json --port 8888
```

## Basic Requests

### Health Check
```bash
curl http://localhost:8888/health
# Expected: {"status":"healthy"}
```

### Get All Users
```bash
curl http://localhost:8888/users
# Expected: [{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]
```

### Get Single User
```bash
curl http://localhost:8888/users/1
# Expected: {"id":1,"name":"Alice","email":"alice@example.com"}
```

### Create User (POST)
```bash
curl -X POST http://localhost:8888/users
# Expected: {"id":3,"name":"Charlie"}
```

### Update User (PUT)
```bash
curl -X PUT http://localhost:8888/users/1
# Expected: {"id":1,"name":"Alice Updated"}
```

### Delete User
```bash
curl -X DELETE http://localhost:8888/users/1 -v
# Expected: 204 No Content
```

## Query Parameter Tests

### Error Simulation
Force a 500 Internal Server Error response:
```bash
curl http://localhost:8888/health?error
# Expected: HTTP 500 with {"status":"healthy"} body
```

### Slow Response
Force a 405 Method Not Allowed response:
```bash
curl http://localhost:8888/health?slow
# Expected: HTTP 405 with {"status":"healthy"} body
```

### Delay Responses

#### 100 Millisecond Delay
```bash
time curl http://localhost:8888/health?delay=100m
# Expected: Response after ~100ms delay
```

#### 2 Second Delay
```bash
time curl http://localhost:8888/health?delay=2s
# Expected: Response after ~2s delay
```

#### 1 Minute Delay
```bash
time curl http://localhost:8888/slow-endpoint?delay=1M
# Expected: Response after ~1 minute delay
```

### Combined Parameters

Delay with error:
```bash
time curl http://localhost:8888/health?delay=1s&error
# Expected: HTTP 500 after ~1s delay
```

Delay with slow:
```bash
time curl http://localhost:8888/users?delay=500m&slow
# Expected: HTTP 405 after ~500ms delay
```

## 404 Not Found

Request a non-existent endpoint:
```bash
curl http://localhost:8888/nonexistent -v
# Expected: HTTP 404 Not Found
```

Request with wrong HTTP verb:
```bash
curl -X PATCH http://localhost:8888/users
# Expected: HTTP 404 Not Found (PATCH not configured)
```

## Headers

Check response headers:
```bash
curl -I http://localhost:8888/health
# Expected: Content-Type: application/json
```

## Verbose Output

Get detailed request/response information:
```bash
curl -v http://localhost:8888/users
```

## Testing with JSON Data

Although WonkyServer doesn't process request bodies (it returns configured responses), you can send data:
```bash
curl -X POST http://localhost:8888/users \
  -H "Content-Type: application/json" \
  -d '{"name":"David","email":"david@example.com"}'
# Expected: {"id":3,"name":"Charlie"} (configured response)
```

## Wonky Mode Testing

Wonky mode introduces controlled chaos for testing unreliable services.

### Setup

Start the server with wonky mode enabled:

```bash
# 50% chance of wonky behavior
./wonkyserver --file example-config.json --wonky 50

# 100% chance (always wonky)
./wonkyserver --file example-config.json --wonky 100

# Combined with custom port
./wonkyserver --file example-config.json --port 9000 --wonky 30
```

### Testing Wonky Behavior

With wonky mode enabled, each request has a chance to randomly exhibit one of three behaviors:

```bash
# Make multiple requests to see different behaviors
curl http://localhost:8888/health  # Might return normally
curl http://localhost:8888/health  # Might return 500 (error)
curl http://localhost:8888/health  # Might delay 5 seconds
curl http://localhost:8888/health  # Might return 405 (slow)
```

### Wonky with Explicit Parameters

Explicit query parameters always override wonky behavior:

```bash
# Start server with 100% wonky
./wonkyserver --file example-config.json --wonky 100

# Explicit params override wonky - these always behave as specified
curl http://localhost:8888/health?error    # Always 500
curl http://localhost:8888/health?slow     # Always 405
curl http://localhost:8888/health?delay=2s # Always delays 2 seconds
```

### Chaos Testing Example

Simulate an unreliable service for testing error handling:

```bash
# Terminal 1: Start server with 40% wonky
./wonkyserver --file example-config.json --wonky 40

# Terminal 2: Test your client's resilience
for i in {1..20}; do
  echo "Request $i:"
  curl -w "\nHTTP Status: %{http_code}\n" http://localhost:8888/users
  echo "---"
done
```

Expected output: Mix of successful responses, 500 errors, 405 responses, and delayed responses.

### Logging

When wonky mode is enabled, the server logs when wonky behavior is applied:

```
2025/01/15 10:30:00 Wonky mode enabled with 50% likelihood
2025/01/15 10:30:05 Received GET request for /health
2025/01/15 10:30:05 Wonky behavior applied: error
2025/01/15 10:30:05 Wonky error: returning 500
2025/01/15 10:30:05 Returned 500 status code
```
