# No-as-a-Service Default Route Plugin for Traefik

This Traefik plugin displays messages from [naas.isalman.dev](https://naas.isalman.dev) as a beautifully styled HTML page with automatic light/dark mode support. The API is called **client-side** via JavaScript to avoid rate limiting on the server. Perfect for a default/fallback route that handles **all unmatched URLs**.

## Features

- ✨ **Beautiful Design** - Responsive layout with smooth animations
- 🌓 **Light/Dark Mode** - Automatic theme switching with manual toggle
- 🔄 **Graceful Fallback** - Configurable default message when API times out
- 🎯 **Catch-All Routes** - Intercepts all requested URLs, not just `/`
- 🚀 **No Extra Containers** - Pure Traefik middleware plugin
- 📱 **Client-Side Fetching** - Each visitor fetches their own message, avoiding server-side rate limits

## Quick Start

```bash
# Clone or download the plugin files
git clone https://github.com/yourusername/noaas-default-route.git
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
      moduleName: github.com/yourusername/noaas-default-route
```

### Option 2: GitHub Plugin (Production)

1. Create a GitHub repository (e.g., `yourusername/noaas-default-route`)

2. Push these files to the repository

3. Configure Traefik to use the plugin in your `traefik.yml`:
```yaml
experimental:
  plugins:
    noaas-default-route:
      moduleName: github.com/yourusername/noaas-default-route
      version: v0.1.0
```

## Configuration

### Static Configuration (traefik.yml)

```yaml
experimental:
  plugins:
    noaas-default-route:
      moduleName: github.com/yourusername/noaas-default-route
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
          apiEndpoint: "https://no-as-a-service.fly.dev/api"
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
      - "--experimental.plugins.noaas-default-route.modulename=github.com/yourusername/noaas-default-route"
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
      - "traefik.http.middlewares.noaas-default.plugin.noaas-default-route.apiEndpoint=https://no-as-a-service.fly.dev/api"
      - "traefik.http.middlewares.noaas-default.plugin.noaas-default-route.defaultMessage=Go Away"
```

## How It Works

1. The plugin intercepts **all** requests to unmatched routes (not just `/`)
2. It generates a beautiful HTML page with embedded JavaScript
3. When the page loads in the visitor's browser, it fetches a message from the API client-side
4. The API returns JSON: `{"reason": "nope"}` 
5. If the API times out or fails, it uses the configured `defaultMessage` instead
6. The page automatically adapts to light/dark mode based on system preferences
7. Visitors can manually toggle between light and dark mode with the button

## Configuration Options

- `apiEndpoint` (optional): The API endpoint for client-side fetching. Defaults to `https://naas.isalman.dev/no`
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

**Important:** Make sure you're running the tests from within the plugin directory!

```bash
# Navigate to the plugin directory first
cd noaas-default-route

# Then run tests
make test

# Or manually
go test -v -race -coverprofile=coverage.out ./...

# Or use the test runner script
./run-tests.sh
```

**Common issue:** If you see `[no test files]`, you're likely running the command from the wrong directory. Make sure you're inside the `noaas-default-route` directory where the `.go` files are located.

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

**Troubleshooting 0% coverage:**

If you see `coverage: 0.0% of statements`, this might be because:

1. **Wrong module name**: Update `go.mod` to match your actual repository path
   ```bash
   # In go.mod, change:
   module github.com/yourusername/noaas-default-route
   # To your actual path:
   module github.com/r-win/noaas-default-route
   ```

2. **Check if tests are actually running**:
   ```bash
   go test -v ./...
   ```
   You should see test output like `TestCreateConfig`, `TestNew`, etc.

3. **Verify coverage file**:
   ```bash
   head coverage.out
   ```
   Should show lines starting with the package name

Expected coverage should be around 80-90% with the provided tests.

## Notes

- The API is called **client-side** by each visitor's browser, not by the Traefik server
- This avoids rate limiting issues on the server side
- The priority should be set to 1 (or low) so it acts as a catch-all route for **all** unmatched URLs
- The HTML is fully responsive and mobile-friendly
- Light/dark mode switches automatically based on system preferences using `prefers-color-scheme`
- Users can manually toggle themes with the button in the top-right corner
- The theme preference is saved in localStorage
- The design includes smooth animations, swaying dot pattern, and hover effects
