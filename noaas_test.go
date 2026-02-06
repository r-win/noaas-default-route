package noaas_default_route

import (
	"context"
	"encoding/json"
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

func TestServeHTTP_Success(t *testing.T) {
	// Mock API server that returns a successful response
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ReasonResponse{Reason: "nope"}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer mockAPI.Close()

	config := &Config{
		APIEndpoint:    mockAPI.URL,
		DefaultMessage: "Go Away",
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
	if !strings.Contains(body, "nope") {
		t.Errorf("Expected body to contain 'nope', got: %s", body)
	}

	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("Expected body to contain HTML doctype")
	}
}

func TestServeHTTP_APIFailure(t *testing.T) {
	// Mock API server that always fails
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockAPI.Close()

	config := &Config{
		APIEndpoint:    mockAPI.URL,
		DefaultMessage: "Custom Fallback",
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

	body := rec.Body.String()
	if !strings.Contains(body, "Custom Fallback") {
		t.Errorf("Expected body to contain 'Custom Fallback', got: %s", body)
	}
}

func TestServeHTTP_APITimeout(t *testing.T) {
	// Mock API server that takes too long to respond
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a long delay (longer than the 5s timeout)
		// Note: In a real test, we'd actually sleep, but for unit tests we just close the connection
		hj, ok := w.(http.Hijacker)
		if ok {
			conn, _, err := hj.Hijack()
			if err != nil {
				t.Errorf("Failed to hijack connection: %v", err)
				return
			}
			_ = conn.Close()
		}
	}))
	defer mockAPI.Close()

	config := &Config{
		APIEndpoint:    mockAPI.URL,
		DefaultMessage: "Timeout Message",
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

	body := rec.Body.String()
	if !strings.Contains(body, "Timeout Message") {
		t.Errorf("Expected body to contain fallback message, got: %s", body)
	}
}

func TestServeHTTP_InvalidJSON(t *testing.T) {
	// Mock API server that returns invalid JSON
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte("invalid json{{{"))
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer mockAPI.Close()

	config := &Config{
		APIEndpoint:    mockAPI.URL,
		DefaultMessage: "JSON Error",
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

	body := rec.Body.String()
	if !strings.Contains(body, "JSON Error") {
		t.Errorf("Expected body to contain fallback message, got: %s", body)
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
		name     string
		message  string
		expected []string
	}{
		{
			name:    "Simple message",
			message: "nope",
			expected: []string{
				"<!DOCTYPE html>",
				"<html lang=\"en\">",
				"<title>No as a Service</title>",
				"nope",
				"no-as-a-service",
				"prefers-color-scheme",
			},
		},
		{
			name:    "Custom message",
			message: "Go Away",
			expected: []string{
				"Go Away",
				"<!DOCTYPE html>",
			},
		},
		{
			name:    "Message with special characters",
			message: "Absolutely not!",
			expected: []string{
				"Absolutely not!",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := plugin.generateHTML(tt.message)

			for _, exp := range tt.expected {
				if !strings.Contains(html, exp) {
					t.Errorf("Expected HTML to contain '%s', but it didn't.\nGenerated HTML:\n%s", exp, html)
				}
			}
		})
	}
}

func TestFetchNoMessage(t *testing.T) {
	tests := []struct {
		name        string
		mockHandler http.HandlerFunc
		wantErr     bool
		wantMessage string
	}{
		{
			name: "Successful fetch",
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				response := ReasonResponse{Reason: "absolutely not"}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(response); err != nil {
					t.Errorf("Failed to encode response: %v", err)
				}
			},
			wantErr:     false,
			wantMessage: "absolutely not",
		},
		{
			name: "Server error",
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "Invalid JSON",
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte("not json"))
				if err != nil {
					t.Errorf("Failed to write response: %v", err)
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := httptest.NewServer(tt.mockHandler)
			defer mockAPI.Close()

			config := &Config{APIEndpoint: mockAPI.URL}
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})
			handler, err := New(context.Background(), next, config, "test")
			if err != nil {
				t.Fatalf("Failed to create handler: %v", err)
			}

			plugin, ok := handler.(*NoaaSDefaultRoute)
			if !ok {
				t.Fatal("Handler is not of type *NoaaSDefaultRoute")
			}

			message, err := plugin.fetchNoMessage()

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.wantErr && message != tt.wantMessage {
				t.Errorf("Expected message '%s', got '%s'", tt.wantMessage, message)
			}
		})
	}
}

func TestServeHTTP_AllURLs(t *testing.T) {
	// Mock API server
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ReasonResponse{Reason: "nope"}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer mockAPI.Close()

	config := &Config{
		APIEndpoint:    mockAPI.URL,
		DefaultMessage: "Go Away",
	}

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t.Fatal("Next handler should not be called")
	})

	handler, err := New(context.Background(), next, config, "test")
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	// Test various URLs - all should be intercepted
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
			if !strings.Contains(body, "nope") {
				t.Errorf("Expected body to contain 'nope' for URL %s", url)
			}
		})
	}
}
