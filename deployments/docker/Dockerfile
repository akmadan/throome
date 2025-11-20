# Multi-stage Dockerfile for Throome Gateway

# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binaries
ARG VERSION=dev
ARG BUILD_TIME
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" \
    -o throome ./cmd/throome && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" \
    -o throome-cli ./cmd/throome-cli

# Stage 2: Runtime
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 throome && \
    adduser -D -u 1000 -G throome throome

# Set working directory
WORKDIR /app

# Copy binaries from builder
COPY --from=builder /build/throome /usr/local/bin/throome
COPY --from=builder /build/throome-cli /usr/local/bin/throome-cli

# Copy configuration examples
COPY --from=builder /build/configs/*.yaml /app/configs/

# Create directories
RUN mkdir -p /app/clusters /app/logs && \
    chown -R throome:throome /app

# Switch to non-root user
USER throome

# Expose ports
EXPOSE 9000 9001

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9000/health || exit 1

# Set environment variables
ENV THROOME_PORT=9000
ENV THROOME_CLUSTERS_DIR=/app/clusters

# Default command
ENTRYPOINT ["throome"]
CMD ["--port", "9000", "--clusters-dir", "/app/clusters"]

