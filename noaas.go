package noaas_default_route

import (
	"context"
	"fmt"
	"net/http"
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

// NoaaSDefaultRoute a Traefik plugin that displays messages from no-as-a-service.
type NoaaSDefaultRoute struct {
	next           http.Handler
	apiEndpoint    string
	defaultMessage string
	name           string
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
	}, nil
}

func (n *NoaaSDefaultRoute) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Generate HTML response with client-side API call
	// Pass the API endpoint and default message to the HTML
	html := n.generateHTML(n.apiEndpoint, n.defaultMessage)

	// Set headers and write response
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write([]byte(html)); err != nil {
		// Log error but don't fail - response already started
		return
	}
}

func (n *NoaaSDefaultRoute) generateHTML(apiEndpoint, defaultMessage string) string {
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

        [data-theme="light"] {
            --bg-gradient-start: #667eea;
            --bg-gradient-end: #764ba2;
            --card-bg: #ffffff;
            --text-primary: #1a202c;
            --text-secondary: #4a5568;
            --link-color: #667eea;
            --shadow-color: rgba(0, 0, 0, 0.1);
            --shadow-hover: rgba(0, 0, 0, 0.15);
        }

        [data-theme="dark"] {
            --bg-gradient-start: #1a1a2e;
            --bg-gradient-end: #16213e;
            --card-bg: #0f1419;
            --text-primary: #e2e8f0;
            --text-secondary: #a0aec0;
            --link-color: #818cf8;
            --shadow-color: rgba(0, 0, 0, 0.3);
            --shadow-hover: rgba(0, 0, 0, 0.5);
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
            transition: background 0.3s ease;
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
            animation: sway 5s ease-in-out infinite;
            pointer-events: none;
        }

        @keyframes sway {
            0%% { transform: translate(0, 0); }
            50%% { transform: translate(25px, 25px); }
            100%% { transform: translate(0, 50px); }
        }

        .theme-toggle {
            position: fixed;
            top: 2rem;
            right: 2rem;
            background: var(--card-bg);
            border: none;
            padding: 0.75rem;
            border-radius: 50%%;
            cursor: pointer;
            box-shadow: 0 4px 12px var(--shadow-color);
            font-size: 1.5rem;
            transition: transform 0.2s ease, box-shadow 0.2s ease;
            z-index: 100;
            width: 50px;
            height: 50px;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .theme-toggle:hover {
            transform: scale(1.1) rotate(15deg);
            box-shadow: 0 6px 16px var(--shadow-hover);
        }

        .theme-toggle:active {
            transform: scale(0.95);
        }

        .container {
            position: relative;
            text-align: center;
            padding: 4rem 3rem;
            background: var(--card-bg);
            border-radius: 2rem;
            box-shadow: 0 25px 70px var(--shadow-color);
            max-width: 85%%;
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

        .loading {
            opacity: 0.5;
        }

        @media (max-width: 640px) {
            .container {
                padding: 3rem 2rem;
            }

            body {
                padding: 1rem;
            }

            .theme-toggle {
                top: 1rem;
                right: 1rem;
            }
        }
    </style>
</head>
<body>
    <button class="theme-toggle" onclick="toggleTheme()" title="Toggle theme" aria-label="Toggle theme">
        <span id="theme-icon">🌙</span>
    </button>
    <div class="container">
        <h1 id="message" class="loading">Loading...</h1>
        <p>Powered by <a href="https://naas.isalman.dev" target="_blank" rel="noopener noreferrer">naas.isalman.dev</a></p>
    </div>

    <script>
        const API_ENDPOINT = '%s';
        const DEFAULT_MESSAGE = '%s';

        // Fetch message from API
        async function fetchMessage() {
            const messageEl = document.getElementById('message');
            
            try {
                const response = await fetch(API_ENDPOINT);
                
                if (!response.ok) {
                    throw new Error('API request failed');
                }
                
                const data = await response.json();
                messageEl.textContent = data.reason || DEFAULT_MESSAGE;
            } catch (error) {
                console.error('Failed to fetch message:', error);
                messageEl.textContent = DEFAULT_MESSAGE;
            } finally {
                messageEl.classList.remove('loading');
            }
        }

        // Theme management
        const getPreferredTheme = () => {
            const savedTheme = localStorage.getItem('theme');
            if (savedTheme) {
                return savedTheme;
            }
            return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
        };

        const setTheme = (theme) => {
            document.documentElement.setAttribute('data-theme', theme);
            localStorage.setItem('theme', theme);
            document.getElementById('theme-icon').textContent = theme === 'dark' ? '☀️' : '🌙';
        };

        const toggleTheme = () => {
            const currentTheme = document.documentElement.getAttribute('data-theme') || getPreferredTheme();
            const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
            setTheme(newTheme);
        };

        // Initialize on page load
        setTheme(getPreferredTheme());
        fetchMessage();

        // Listen for system theme changes
        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
            if (!localStorage.getItem('theme')) {
                setTheme(e.matches ? 'dark' : 'light');
            }
        });
    </script>
</body>
</html>`, apiEndpoint, defaultMessage)
}
