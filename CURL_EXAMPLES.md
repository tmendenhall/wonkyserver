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
