# No-as-a-Service Default Route Plugin for Traefik

This Traefik plugin fetches a message from [no-as-a-service](https://naas.isalman.dev) and returns it as a beautifully styled HTML page with automatic light/dark mode support. Perfect for a default/fallback route that handles **all unmatched URLs**.

## Features

- ✨ **Beautiful Design** - Responsive layout with smooth animations
- 🌓 **Light/Dark Mode** - Automatically adapts to system preferences
- 🔄 **Graceful Fallback** - Configurable default message when API times out
- 🎯 **Catch-All Routes** - Intercepts all requested URLs, not just `/`
- 🚀 **No Extra Containers** - Pure Traefik middleware plugin

## Quick Start

```bash
# Clone or download the plugin files
git clone https://github.com/r-win/noaas-default-route.git
cd noaas-default-route

# Run the setup script (installs golangci-lint if needed)
./setup.sh

# Or install dependencies manually
make install-lint
go mod download

# Run tests
make test
```

## Installation

### Option 1: Local Plugin (Development)

1. Create the plugin directory structure:
```bash
mkdir -p /path/to/plugins/noaas-default-route
```

2. Copy the plugin files (`noaas.go` and `.traefik.yml`) to this directory.

3. Configure Traefik to use local plugins in your `traefik.yml` or static configuration:
```yaml
experimental:
  localPlugins:
    noaas-default-route:
      moduleName: github.com/r-win/noaas-default-route
```

### Option 2: GitHub Plugin (Production)

1. Create a GitHub repository (e.g., `r-win/noaas-default-route`)

2. Push these files to the repository

3. Configure Traefik to use the plugin in your `traefik.yml`:
```yaml
experimental:
  plugins:
    noaas-default-route:
      moduleName: github.com/r-win/noaas-default-route
      version: v0.1.0
```

## Configuration

### Static Configuration (traefik.yml)

```yaml
experimental:
  plugins:
    noaas-default-route:
      moduleName: github.com/r-win/noaas-default-route
      version: v0.1.0

entryPoints:
  web:
    address: ":80"
```

### Dynamic Configuration

Create a file like `dynamic-config.yml`:

```yaml
http:
  routers:
    default-router:
      rule: "PathPrefix(`/`)"
      service: noaas-service
      middlewares:
        - noaas-default
      priority: 1  # Lowest priority so it acts as catch-all

  middlewares:
    noaas-default:
      plugin:
        noaas-default-route:
          apiEndpoint: "https://naas.isalman.dev/no"
          defaultMessage: "Go Away"

  services:
    noaas-service:
      loadBalancer:
        servers:
          - url: "http://localhost:9999"  # Dummy backend (won't be reached)
```

Or use Docker labels if you're using Docker:

```yaml
services:
  traefik:
    image: traefik:v3.0
    command:
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--entrypoints.web.address=:80"
      - "--experimental.plugins.noaas-default-route.modulename=github.com/r-win/noaas-default-route"
      - "--experimental.plugins.noaas-default-route.version=v0.1.0"
    ports:
      - "80:80"
      - "8080:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock

  # Dummy service for the default route
  default-route:
    image: traefik/whoami
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.default.rule=PathPrefix(`/`)"
      - "traefik.http.routers.default.priority=1"
      - "traefik.http.routers.default.middlewares=noaas-default@docker"
      - "traefik.http.middlewares.noaas-default.plugin.noaas-default-route.apiEndpoint=https://naas.isalman.dev/no"
      - "traefik.http.middlewares.noaas-default.plugin.noaas-default-route.defaultMessage=Go Away"
```

## How It Works

1. The plugin intercepts **all** requests to unmatched routes (not just `/`)
2. It makes an HTTP GET request to the no-as-a-service API
3. It receives a JSON response like `{"no": "nope"}`
4. If the API times out or fails, it uses the configured `defaultMessage` instead
5. It generates a beautiful, responsive HTML page with the message
6. The page automatically adapts to light/dark mode based on system preferences
7. It returns the HTML to the client

## Configuration Options

- `apiEndpoint` (optional): The API endpoint to fetch messages from. Defaults to `https://naas.isalman.dev/no`
- `defaultMessage` (optional): Message to display when API times out or fails. Defaults to `Go Away`

## Example

When you visit your Traefik instance at the default route, you'll see a beautifully styled page displaying one of the many variations of "no" from the no-as-a-service API.

## Testing Locally

You can test the plugin with Traefik's local plugin support:

```bash
# Start Traefik with the configuration
traefik --configFile=traefik.yml

# Visit http://localhost to see the default route
curl http://localhost
```

## Development

### Prerequisites

- Go 1.21 or later
- golangci-lint (for linting)

### Setup

```bash
# Install golangci-lint
make install-lint

# Or install manually
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest

# Optional: Add Go bin to PATH if not already added
# Add this to your ~/.bashrc or ~/.zshrc:
export PATH="$PATH:$(go env GOPATH)/bin"

# Verify installation
$(go env GOPATH)/bin/golangci-lint --version
# Or if $(go env GOPATH)/bin is in your PATH:
golangci-lint --version
```

**Note:** The Makefile automatically uses the full path to golangci-lint, so you don't need to add it to your PATH for `make` commands to work. However, adding it to PATH is useful for running `golangci-lint` directly.

### Running Tests

```bash
# Run all tests
make test

# Or manually
go test -v -race -coverprofile=coverage.out ./...
```

### Linting

```bash
# Run golangci-lint
make lint

# Or manually
golangci-lint run
```

### Code Formatting

```bash
# Format code
make fmt

# Run all checks (fmt, vet, lint, test)
make check
```

### Coverage

After running tests, view the coverage report:
```bash
go tool cover -html=coverage.out
```

## Notes

- The plugin has a 5-second timeout for API requests
- If the API request fails, it uses the configured `defaultMessage` (defaults to "Go Away")
- The priority should be set to 1 (or low) so it acts as a catch-all route for **all** unmatched URLs
- The HTML is fully responsive and mobile-friendly
- Light/dark mode switches automatically based on system preferences using `prefers-color-scheme`
- The design includes smooth animations and hover effects
