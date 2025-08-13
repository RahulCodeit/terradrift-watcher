# Multi-stage build for TerraDrift Watcher
# This creates a minimal container image with the tool and Terraform

# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# -ldflags "-s -w" strips debug info for smaller size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o terradrift-watcher .

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    terraform \
    git \
    bash \
    curl \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1000 terradrift && \
    adduser -D -u 1000 -G terradrift terradrift

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/terradrift-watcher /app/terradrift-watcher
COPY --from=builder /build/config.example.yml /app/config.example.yml

# Make binary executable
RUN chmod +x /app/terradrift-watcher

# Create directories for configs and terraform projects
RUN mkdir -p /config /terraform && \
    chown -R terradrift:terradrift /app /config /terraform

# Switch to non-root user
USER terradrift

# Set environment variables
ENV PATH="/app:${PATH}"

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD terradrift-watcher --version || exit 1

# Default command
ENTRYPOINT ["/app/terradrift-watcher"]
CMD ["--help"] 