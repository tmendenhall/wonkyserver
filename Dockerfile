# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files (will be created when user runs go mod init)
COPY go.mod* ./

# Download dependencies if go.mod exists
RUN if [ -f go.mod ]; then go mod download; fi

# Copy source code
COPY *.go ./

# Build the application
# CGO_ENABLED=0 for static binary
# -ldflags="-w -s" to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o wonkyserver .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/wonkyserver .

# Expose default port
EXPOSE 8888

# Run as non-root user
RUN adduser -D -u 1000 appuser
USER appuser

# Default command
ENTRYPOINT ["/app/wonkyserver"]
