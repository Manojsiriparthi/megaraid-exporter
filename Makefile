.PHONY: build test clean install test-megacli production-check

# Variables
BINARY_NAME=megaraid-exporter
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Default target
all: build

# Build the binary
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/exporter

# Run tests
test:
	go test -v ./...

# Run integration tests with MegaCLI (requires hardware)
test-integration:
	sudo go test -tags=integration -v ./...

# Test MegaCLI availability
test-megacli:
	@echo "Testing MegaCLI64 availability..."
	@which megacli64 >/dev/null 2>&1 || (echo "ERROR: megacli64 not found in PATH" && exit 1)
	@sudo megacli64 -v || (echo "ERROR: Cannot execute megacli64" && exit 1)
	@echo "MegaCLI64 is available and working"

# Test HTTP endpoint
test-endpoint:
	@echo "Testing HTTP endpoint..."
	@curl -s http://localhost:9272/metrics > /dev/null || (echo "ERROR: Exporter not responding on port 9272" && exit 1)
	@echo "HTTP endpoint is responding"

# Production readiness check
production-check: test-megacli
	@echo "=== Production Readiness Check ==="
	@echo "Checking system requirements..."
	@which systemctl >/dev/null 2>&1 || (echo "ERROR: systemctl not found (systemd required)" && exit 1)
	@which curl >/dev/null 2>&1 || (echo "ERROR: curl not found" && exit 1)
	@[ $$(id -u) -eq 0 ] || (echo "ERROR: Must run as root for production deployment" && exit 1)
	@echo "✓ All production requirements met"

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	go clean
	rm -rf dist/

# Install binary to system (production)
install-production: build production-check
	@echo "Installing for production..."
	sudo cp $(BINARY_NAME) /usr/local/bin/
	sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	sudo mkdir -p /etc/$(BINARY_NAME)
	sudo mkdir -p /var/log/$(BINARY_NAME)
	@echo "✓ Production installation completed"

# Cross-compile for different platforms
build-all:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 ./cmd/exporter
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 ./cmd/exporter
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe ./cmd/exporter

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Development setup check
dev-check: test-megacli
	@echo "Checking Go version..."
	@go version
	@echo "Checking dependencies..."
	@go mod verify
	@echo "Development environment ready"

# Quick smoke test
smoke-test: build
	@echo "Running smoke test..."
	@timeout 10s sudo ./$(BINARY_NAME) --log-level debug &
	@sleep 3
	@curl -s http://localhost:9272/metrics | head -5
	@pkill -f $(BINARY_NAME) || true
	@echo "Smoke test completed"

# Show help
help:
	@echo "Available targets:"
	@echo "  build              - Build the binary"
	@echo "  test               - Run unit tests"
	@echo "  test-integration   - Run integration tests (requires MegaCLI + hardware)"
	@echo "  test-megacli       - Test MegaCLI64 availability"
	@echo "  test-endpoint      - Test HTTP endpoint availability"
	@echo "  production-check   - Verify production requirements"
	@echo "  install-production - Install for production use"
	@echo "  clean              - Clean build artifacts"
	@echo "  build-all          - Cross-compile for multiple platforms"
	@echo "  fmt                - Format code"
	@echo "  lint               - Run linter"
	@echo "  dev-check          - Check development environment"
	@echo "  smoke-test         - Quick build and endpoint test"
