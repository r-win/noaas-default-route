.PHONY: test lint fmt vet build clean help install-lint check-simple

# Get golangci-lint path
GOPATH ?= $(shell go env GOPATH)
GOLANGCI_LINT := $(GOPATH)/bin/golangci-lint

# Default target
help:
	@echo "Available targets:"
	@echo "  test         - Run tests"
	@echo "  lint         - Run golangci-lint"
	@echo "  install-lint - Install golangci-lint"
	@echo "  fmt          - Format code with gofmt"
	@echo "  vet          - Run go vet"
	@echo "  build        - Build the plugin"
	@echo "  clean        - Clean build artifacts"
	@echo "  check        - Run fmt, vet, lint, and test"
	@echo "  check-simple - Run fmt, vet, and test (no linter)"

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Get golangci-lint path
GOPATH ?= $(shell go env GOPATH)
GOLANGCI_LINT := $(GOPATH)/bin/golangci-lint

# Install golangci-lint
install-lint:
	@echo "Installing golangci-lint..."
	@test -f $(GOLANGCI_LINT) || \
		(curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin latest)
	@echo "golangci-lint installed successfully at $(GOLANGCI_LINT)"

# Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@test -f $(GOLANGCI_LINT) || \
		(echo "golangci-lint not found at $(GOLANGCI_LINT). Run 'make install-lint' to install it." && exit 1)
	$(GOLANGCI_LINT) run

# Format code
fmt:
	@echo "Formatting code..."
	gofmt -s -w .

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Build the plugin (verify it compiles)
build:
	@echo "Building plugin..."
	go build -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f coverage.out coverage.html

# Run all checks
check: fmt vet lint test
	@echo "All checks passed!"

# Run checks without linter (for quick local testing)
check-simple: fmt vet test
	@echo "Simple checks passed! (Run 'make check' for full linting)"
