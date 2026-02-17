# Multi-stage Dockerfile for Audoctl
# Stage 1: Build
FROM golang:1.26-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build arguments for version information
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X github.com/audoctl/audoctl/cmd/version.Version=${VERSION} -X github.com/audoctl/audoctl/cmd/version.Commit=${COMMIT} -X github.com/audoctl/audoctl/cmd/version.BuildTime=${BUILD_TIME}" \
    -o /app/bin/audoctl \
    ./cmd

# Stage 2: Runtime
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 audoctl && \
    adduser -D -u 1000 -G audoctl audoctl

# Set working directory
WORKDIR /home/audoctl

# Copy binary from builder
COPY --from=builder /app/bin/audoctl /usr/local/bin/audoctl

# Change ownership
RUN chown -R audoctl:audoctl /home/audoctl

# Switch to non-root user
USER audoctl

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
ENTRYPOINT ["/usr/local/bin/audoctl"]
