# Audoctl Makefile
.PHONY: help build run clean test lint fmt install dev

# Variables
BINARY_NAME=audoctl
VERSION?=dev
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X github.com/audoctl/audoctl/cmd/version.Version=$(VERSION) -X github.com/audoctl/audoctl/cmd/version.Commit=$(COMMIT) -X github.com/audoctl/audoctl/cmd/version.BuildTime=$(BUILD_TIME)"

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd
	@echo "✓ Build complete: bin/$(BINARY_NAME)"

run: build ## Build and run the application
	@echo "Starting $(BINARY_NAME)..."
	@./bin/$(BINARY_NAME) audoctl

dev: ## Run in development mode with hot reload (requires air)
	@if ! command -v air > /dev/null; then \
		echo "Installing air..."; \
		go install github.com/air-verse/air@latest; \
	fi
	@air

install: ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) ./cmd
	@echo "✓ Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

clean: ## Remove built binaries and temporary files
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f audoctl.db
	@go clean
	@echo "✓ Clean complete"

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "✓ Tests complete"

test-coverage: test ## Run tests with coverage report
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

lint: ## Run linter (requires golangci-lint)
	@if ! command -v golangci-lint > /dev/null; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo "Running linter..."
	@golangci-lint run ./...
	@echo "✓ Lint complete"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -s -w .
	@echo "✓ Format complete"

tidy: ## Tidy go modules
	@echo "Tidying modules..."
	@go mod tidy
	@echo "✓ Tidy complete"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ Vet complete"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@echo "✓ Dependencies downloaded"

check: fmt vet lint test ## Run all checks (format, vet, lint, test)
	@echo "✓ All checks passed"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):$(VERSION) .
	@echo "✓ Docker image built: $(BINARY_NAME):$(VERSION)"

version: ## Display version information
	@./bin/$(BINARY_NAME) version || echo "Build the binary first with 'make build'"

migrate-up: ## Run database migrations (up)
	@echo "Running migrations up..."
	@./bin/$(BINARY_NAME) migrate up

migrate-down: ## Run database migrations (down)
	@echo "Running migrations down..."
	@./bin/$(BINARY_NAME) migrate down

generate:
	swagger generate spec ./cmd/audoctl/docs/ -o ./cmd/audoctl/docs/swagger.json --scan-models

.DEFAULT_GOAL := help
