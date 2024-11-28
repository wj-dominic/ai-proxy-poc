package loadbalancer

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/wj-dominic/ai-proxy-poc/pkg/node"
)

type mockAlgorithm struct {
	nodes []*node.Node
	index int
}

func (m *mockAlgorithm) NextNode() *node.Node {
	if len(m.nodes) == 0 {
		return nil
	}
	node := m.nodes[m.index]
	m.index = (m.index + 1) % len(m.nodes)
	return node
}

func TestLoadBalancer_ServeHTTP_NoNodes(t *testing.T) {
	lb := NewLoadBalancer(nil, &mockAlgorithm{nodes: nil})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(lb.ServeHTTP)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, rr.Code)
	}

	expectedBody := "No nodes available\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}
}

func TestLoadBalancer_ServeHTTP_WithNodes(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from Backend"))
	}))
	defer backend.Close()

	backendURL, err := url.Parse(backend.URL)
	if err != nil {
		t.Fatalf("Failed to parse backend URL: %v", err)
	}

	nodes := []*node.Node{
		node.NewNode("node1", backendURL, 1000, 100),
		node.NewNode("node2", backendURL, 1000, 100),
	}
	lb := NewLoadBalancer(nodes, &mockAlgorithm{nodes: nodes})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(lb.ServeHTTP)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	expectedBody := "Hello from Backend"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}
}
