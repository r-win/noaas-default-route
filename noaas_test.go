package noaas_default_route

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateConfig(t *testing.T) {
	config := CreateConfig()

	if config.APIEndpoint != "https://naas.isalman.dev/no" {
		t.Errorf("Expected default APIEndpoint to be 'https://naas.isalman.dev/no', got '%s'", config.APIEndpoint)
	}

	if config.DefaultMessage != "Go Away" {
		t.Errorf("Expected default DefaultMessage to be 'Go Away', got '%s'", config.DefaultMessage)
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectedAPI string
		expectedMsg string
	}{
		{
			name:        "Default configuration",
			config:      &Config{},
			expectedAPI: "https://naas.isalman.dev/no",
			expectedMsg: "Go Away",
		},
		{
			name: "Custom configuration",
			config: &Config{
				APIEndpoint:    "https://custom-api.example.com/api",
				DefaultMessage: "Custom Message",
			},
			expectedAPI: "https://custom-api.example.com/api",
			expectedMsg: "Custom Message",
		},
		{
			name: "Partial configuration",
			config: &Config{
				APIEndpoint: "https://custom-api.example.com/api",
			},
			expectedAPI: "https://custom-api.example.com/api",
			expectedMsg: "Go Away",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})
			handler, err := New(context.Background(), next, tt.config, "test")

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			plugin, ok := handler.(*NoaaSDefaultRoute)
			if !ok {
				t.Fatal("Handler is not of type *NoaaSDefaultRoute")
			}

			if plugin.apiEndpoint != tt.expectedAPI {
				t.Errorf("Expected apiEndpoint '%s', got '%s'", tt.expectedAPI, plugin.apiEndpoint)
			}

			if plugin.defaultMessage != tt.expectedMsg {
				t.Errorf("Expected defaultMessage '%s', got '%s'", tt.expectedMsg, plugin.defaultMessage)
			}
		})
	}
}

func TestServeHTTP_HTMLGeneration(t *testing.T) {
	config := &Config{
		APIEndpoint:    "https://naas.isalman.dev/no",
		DefaultMessage: "Test Message",
	}

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t.Fatal("Next handler should not be called")
	})

	handler, err := New(context.Background(), next, config, "test")
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/html; charset=utf-8', got '%s'", contentType)
	}

	body := rec.Body.String()

	// Check for HTML structure
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("Expected body to contain HTML doctype")
	}

	if !strings.Contains(body, "<html lang=\"en\">") {
		t.Error("Expected body to contain html tag")
	}

	// Check for embedded API endpoint in JavaScript
	if !strings.Contains(body, "https://naas.isalman.dev/no") {
		t.Error("Expected body to contain API endpoint in JavaScript")
	}

	// Check for default message in JavaScript
	if !strings.Contains(body, "Test Message") {
		t.Error("Expected body to contain default message in JavaScript")
	}

	// Check for loading state
	if !strings.Contains(body, "Loading...") {
		t.Error("Expected body to contain loading message")
	}

	// Check for theme toggle functionality
	if !strings.Contains(body, "toggleTheme") {
		t.Error("Expected body to contain theme toggle function")
	}

	if !strings.Contains(body, "theme-toggle") {
		t.Error("Expected body to contain theme toggle button")
	}

	// Check for fetchMessage function
	if !strings.Contains(body, "fetchMessage") {
		t.Error("Expected body to contain fetchMessage function")
	}

	// Check for proper JSON field access
	if !strings.Contains(body, "data.reason") {
		t.Error("Expected body to access data.reason from API response")
	}
}

func TestServeHTTP_CustomEndpoint(t *testing.T) {
	customEndpoint := "https://example.com/custom/api"
	config := &Config{
		APIEndpoint:    customEndpoint,
		DefaultMessage: "Custom Default",
	}

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t.Fatal("Next handler should not be called")
	})

	handler, err := New(context.Background(), next, config, "test")
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/test-path", http.NoBody)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	body := rec.Body.String()

	// Check custom endpoint is embedded
	if !strings.Contains(body, customEndpoint) {
		t.Errorf("Expected body to contain custom endpoint '%s'", customEndpoint)
	}

	// Check custom default message is embedded
	if !strings.Contains(body, "Custom Default") {
		t.Error("Expected body to contain custom default message")
	}
}

func TestGenerateHTML(t *testing.T) {
	config := CreateConfig()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})
	handler, err := New(context.Background(), next, config, "test")
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	plugin, ok := handler.(*NoaaSDefaultRoute)
	if !ok {
		t.Fatal("Handler is not of type *NoaaSDefaultRoute")
	}

	tests := []struct {
		name        string
		apiEndpoint string
		defaultMsg  string
		expected    []string
	}{
		{
			name:        "Default values",
			apiEndpoint: "https://naas.isalman.dev/no",
			defaultMsg:  "Go Away",
			expected: []string{
				"<!DOCTYPE html>",
				"<html lang=\"en\">",
				"<title>No as a Service</title>",
				"https://naas.isalman.dev/no",
				"Go Away",
				"naas.isalman.dev",
				"@keyframes sway",
				"data.reason",
			},
		},
		{
			name:        "Custom values",
			apiEndpoint: "https://example.com/api",
			defaultMsg:  "Custom Message",
			expected: []string{
				"https://example.com/api",
				"Custom Message",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := plugin.generateHTML(tt.apiEndpoint, tt.defaultMsg)

			for _, exp := range tt.expected {
				if !strings.Contains(html, exp) {
					t.Errorf("Expected HTML to contain '%s', but it didn't", exp)
				}
			}
		})
	}
}

func TestServeHTTP_AllURLs(t *testing.T) {
	config := &Config{
		APIEndpoint:    "https://naas.isalman.dev/no",
		DefaultMessage: "Go Away",
	}

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t.Fatal("Next handler should not be called")
	})

	handler, err := New(context.Background(), next, config, "test")
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	// Test various URLs - all should be intercepted and return HTML
	urls := []string{
		"/",
		"/test",
		"/path/to/resource",
		"/api/v1/endpoint",
		"/some/deep/nested/path",
	}

	for _, url := range urls {
		t.Run("URL: "+url, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, url, http.NoBody)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("Expected status code %d for URL %s, got %d", http.StatusOK, url, rec.Code)
			}

			body := rec.Body.String()
			if !strings.Contains(body, "<!DOCTYPE html>") {
				t.Errorf("Expected HTML response for URL %s", url)
			}

			if !strings.Contains(body, "fetchMessage") {
				t.Errorf("Expected JavaScript fetchMessage function for URL %s", url)
			}
		})
	}
}

func TestServeHTTP_DifferentMethods(t *testing.T) {
	config := CreateConfig()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t.Fatal("Next handler should not be called")
	})

	handler, err := New(context.Background(), next, config, "test")
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodHead,
	}

	for _, method := range methods {
		t.Run("Method: "+method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/", http.NoBody)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("Expected status code %d for method %s, got %d", http.StatusOK, method, rec.Code)
			}
		})
	}
}
