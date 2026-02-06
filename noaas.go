package noaas_default_route

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Config the plugin configuration.
type Config struct {
	APIEndpoint    string `json:"apiEndpoint,omitempty"`
	DefaultMessage string `json:"defaultMessage,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		APIEndpoint:    "https://naas.isalman.dev/no",
		DefaultMessage: "Go Away",
	}
}

// NoaaSDefaultRoute a Traefik plugin that fetches messages from no-as-a-service.
type NoaaSDefaultRoute struct {
	next           http.Handler
	apiEndpoint    string
	defaultMessage string
	name           string
	client         *http.Client
}

// ReasonResponse represents the response from no-as-a-service API.
type ReasonResponse struct {
	Reason string `json:"reason"`
}

// New created a new NoaaSDefaultRoute plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.APIEndpoint == "" {
		config.APIEndpoint = "https://naas.isalman.dev/no"
	}

	if config.DefaultMessage == "" {
		config.DefaultMessage = "Go Away"
	}

	return &NoaaSDefaultRoute{
		next:           next,
		apiEndpoint:    config.APIEndpoint,
		defaultMessage: config.DefaultMessage,
		name:           name,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}, nil
}

func (n *NoaaSDefaultRoute) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Fetch the message from no-as-a-service
	message, err := n.fetchNoMessage()
	if err != nil {
		// Use default message if API fails
		message = n.defaultMessage
	}

	// Generate HTML response
	html := n.generateHTML(message)

	// Set headers and write response
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write([]byte(html)); err != nil {
		// Log error but don't fail - response already started
		return
	}
}

func (n *NoaaSDefaultRoute) fetchNoMessage() (string, error) {
	resp, err := n.client.Get(n.apiEndpoint)
	if err != nil {
		return "", fmt.Errorf("failed to fetch from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var reasonResp ReasonResponse
	if err := json.Unmarshal(body, &reasonResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return reasonResp.Reason, nil
}

func (n *NoaaSDefaultRoute) generateHTML(message string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>No as a Service</title>
    <style>
        :root {
            --bg-gradient-start: #667eea;
            --bg-gradient-end: #764ba2;
            --card-bg: #ffffff;
            --text-primary: #1a202c;
            --text-secondary: #4a5568;
            --link-color: #667eea;
            --shadow-color: rgba(0, 0, 0, 0.1);
            --shadow-hover: rgba(0, 0, 0, 0.15);
        }

        @media (prefers-color-scheme: dark) {
            :root {
                --bg-gradient-start: #1a1a2e;
                --bg-gradient-end: #16213e;
                --card-bg: #0f1419;
                --text-primary: #e2e8f0;
                --text-secondary: #a0aec0;
                --link-color: #818cf8;
                --shadow-color: rgba(0, 0, 0, 0.3);
                --shadow-hover: rgba(0, 0, 0, 0.5);
            }
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Helvetica Neue', sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            background: linear-gradient(135deg, var(--bg-gradient-start) 0%%, var(--bg-gradient-end) 100%%);
            padding: 2rem;
            position: relative;
            overflow: hidden;
        }

        body::before {
            content: '';
            position: absolute;
            top: -50%%;
            left: -50%%;
            width: 200%%;
            height: 200%%;
            background: radial-gradient(circle, rgba(255,255,255,0.1) 1px, transparent 1px);
            background-size: 50px 50px;
            animation: drift 60s linear infinite;
            pointer-events: none;
        }

        @keyframes drift {
            0%% { transform: translate(0, 0); }
            100%% { transform: translate(50px, 50px); }
        }

        .container {
            position: relative;
            text-align: center;
            padding: 4rem 3rem;
            background: var(--card-bg);
            border-radius: 2rem;
            box-shadow: 0 25px 70px var(--shadow-color);
            max-width: 900px;
            width: 100%%;
            backdrop-filter: blur(10px);
            transition: transform 0.3s ease, box-shadow 0.3s ease;
        }

        .container:hover {
            transform: translateY(-5px);
            box-shadow: 0 30px 80px var(--shadow-hover);
        }

        h1 {
            font-size: clamp(3rem, 8vw, 6rem);
            margin: 0;
            color: var(--text-primary);
            font-weight: 800;
            letter-spacing: -0.02em;
            line-height: 1.1;
            background: linear-gradient(135deg, var(--bg-gradient-start), var(--bg-gradient-end));
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            animation: fadeIn 0.8s ease-out;
        }

        @keyframes fadeIn {
            from {
                opacity: 0;
                transform: translateY(20px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        p {
            font-size: 1.125rem;
            color: var(--text-secondary);
            margin-top: 2rem;
            font-weight: 500;
            animation: fadeIn 1s ease-out 0.2s backwards;
        }

        a {
            color: var(--link-color);
            text-decoration: none;
            font-weight: 600;
            transition: all 0.2s ease;
            position: relative;
        }

        a::after {
            content: '';
            position: absolute;
            bottom: -2px;
            left: 0;
            width: 0;
            height: 2px;
            background: var(--link-color);
            transition: width 0.3s ease;
        }

        a:hover::after {
            width: 100%%;
        }

        a:hover {
            opacity: 0.8;
        }

        @media (max-width: 640px) {
            .container {
                padding: 3rem 2rem;
            }

            body {
                padding: 1rem;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>%s</h1>
        <p>Powered by <a href="https://naas.isalman.dev" target="_blank" rel="noopener noreferrer">no-as-a-service</a></p>
    </div>
</body>
</html>`, message)
}
