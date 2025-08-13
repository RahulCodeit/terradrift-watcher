# Makefile for TerraDrift Watcher

# Variables
BINARY_NAME=terradrift-watcher
VERSION?=dev
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X 'github.com/terradrift-watcher/cmd.version=$(VERSION)' -X 'github.com/terradrift-watcher/cmd.commit=$(COMMIT)' -X 'github.com/terradrift-watcher/cmd.date=$(DATE)'"

# Default target
.PHONY: all
all: build

# Build for current platform
.PHONY: build
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -cover -coverprofile=coverage.txt ./...
	go tool cover -html=coverage.txt -o coverage.html

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	rm -f coverage.txt coverage.html
	rm -rf dist/

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
.PHONY: lint
lint:
	golangci-lint run

# Build for all platforms
.PHONY: build-all
build-all: clean
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .

# Install locally
.PHONY: install
install: build
	sudo mv $(BINARY_NAME) /usr/local/bin/

# Run the tool
.PHONY: run
run: build
	./$(BINARY_NAME) run --config config.yml

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make build       - Build for current platform"
	@echo "  make test        - Run tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Run linter"
	@echo "  make build-all   - Build for all platforms"
	@echo "  make install     - Install locally to /usr/local/bin"
	@echo "  make run         - Build and run with config.yml" 