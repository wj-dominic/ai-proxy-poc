package loadbalancer

import (
	"sync"
	"testing"

	"github.com/wj-dominic/ai-proxy-poc/pkg/node"
)

func TestRoundRobin_NextNode(t *testing.T) {
	nodes := []*node.Node{
		node.NewNode("node1", nil, 0, 0),
		node.NewNode("node2", nil, 0, 0),
		node.NewNode("node3", nil, 0, 0),
	}

	rr := NewRoundRobin(nodes)

	tests := []struct {
		expectedNodeID string
	}{
		{"node2"},
		{"node3"},
		{"node1"},
		{"node2"},
	}

	for i, tt := range tests {
		node := rr.NextNode()
		if node.ID() != tt.expectedNodeID {
			t.Errorf("test %d: expected node ID %s, got %s", i, tt.expectedNodeID, node.ID())
		}
	}
}

func TestRoundRobin_NextNodeWithGoroutine(t *testing.T) {
	nodes := []*node.Node{
		node.NewNode("node1", nil, 0, 0),
		node.NewNode("node2", nil, 0, 0),
		node.NewNode("node3", nil, 0, 0),
	}

	rr := NewRoundRobin(nodes)

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			_ = rr.NextNode()
			wg.Done()
		}()
	}
	wg.Wait()

	node := rr.NextNode()
	if node.ID() != "node3" {
		t.Errorf("expected node ID node3, got %s", node.ID())
	}
}
