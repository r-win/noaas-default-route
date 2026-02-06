#!/bin/bash

set -e

echo "🚀 Setting up development environment for noaas-default-route..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    echo "   Visit: https://golang.org/doc/install"
    exit 1
fi

echo "✅ Go $(go version | awk '{print $3}') found"

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    echo "📦 Installing golangci-lint..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest
    echo "✅ golangci-lint installed"
else
    echo "✅ golangci-lint $(golangci-lint --version | head -n1) found"
fi

# Download Go dependencies
echo "📦 Downloading Go dependencies..."
go mod download
echo "✅ Dependencies downloaded"

# Run tests to verify setup
echo "🧪 Running tests to verify setup..."
if go test ./... > /dev/null 2>&1; then
    echo "✅ Tests passed"
else
    echo "⚠️  Some tests failed, but setup is complete"
fi

echo ""
echo "✨ Setup complete! You can now run:"
echo "   make test       - Run tests"
echo "   make lint       - Run linter"
echo "   make check      - Run all checks"
echo "   make help       - Show all available commands"
