package node

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// Mock backend server to simulate the origin server
func mockBackend() (*httptest.Server, *url.URL) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Backend-Header", "BackendValue")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from Backend"))
	}))

	backendURL, _ := url.Parse(backend.URL)
	return backend, backendURL
}

// Test Node's ServeHTTP function
func TestNode_ServeHTTP(t *testing.T) {
	// Step 1: Create a mock backend server
	backend, backendURL := mockBackend()
	defer backend.Close()

	// Step 2: Create a Node instance pointing to the mock backend
	node := NewNode("node1", backendURL, 1000, 100)

	// Step 3: Create a test HTTP request
	req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
	rec := httptest.NewRecorder()

	// Step 4: Call Node's ServeHTTP method
	node.ServeHTTP(rec, req)

	// Step 5: Validate the response from the proxy
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}

	if res.Header.Get("X-Backend-Header") != "BackendValue" {
		t.Errorf("expected X-Backend-Header to be 'BackendValue', got '%s'", res.Header.Get("X-Backend-Header"))
	}
}

// Test Node's Director function for header modification
func TestNode_Director(t *testing.T) {
	// Step 1: Create a mock backend server
	backend, backendURL := mockBackend()
	defer backend.Close()

	// Step 2: Create a Node instance with the custom director function
	node := NewNode("node1", backendURL, 1000, 100)

	// Step 3: Create a test HTTP request
	req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)

	// Use the Director function directly to modify the request
	node.proxy.Director(req)

	if req.Header.Get("X-Forwarded-Host") != "example.com" {
		t.Errorf("expected X-Forwarded-Host to be 'example.com', got '%s'", req.Header.Get("X-Forwarded-Host"))
	}

	if req.Header.Get("X-Origin-Host") != backendURL.Host {
		t.Errorf("expected X-Origin-Host to be '%s', got '%s'", backendURL.Host, req.Header.Get("X-Origin-Host"))
	}

	if req.Host != backendURL.Host {
		t.Errorf("expected Host to be '%s', got '%s'", backendURL.Host, req.Host)
	}
}

// Test Node's IsAllowRequest function
func TestNode_IsAllowRequest(t *testing.T) {
	backend, backendURL := mockBackend()
	defer backend.Close()

	node := NewNode("node1", backendURL, 10, 1000)

	tests := []struct {
		name     string
		bodySize int64
		expected bool
	}{
		{"AllowRequestWithinLimits", 100, true},
		{"DenyRequestExceedingRPM", 100, true}, // This will be denied in the next request
		{"DenyRequestExceedingBPM", 1000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed := node.IsAllowRequest(tt.bodySize)
			if allowed != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, allowed)
			}
		})
	}

	// Test exceeding RPM
	for i := 0; i < 10; i++ {
		node.IsAllowRequest(50)
	}
	if node.IsAllowRequest(50) {
		t.Errorf("expected false, got true")
	}

	// Test exceeding BPM
	node.ResetRateLimit()
	if node.IsAllowRequest(1001) {
		t.Errorf("expected false, got true")
	}
}
