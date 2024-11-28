package loadbalancer

import (
	"sync"

	"github.com/wj-dominic/ai-proxy-poc/pkg/node"
)

type RoundRobin struct {
	nodes        []*node.Node
	currentIndex int

	mu sync.Mutex
}

func NewRoundRobin(nodes []*node.Node) *RoundRobin {
	return &RoundRobin{
		nodes:        nodes,
		currentIndex: 0,
		mu:           sync.Mutex{},
	}
}

func (rr *RoundRobin) NextNode() *node.Node {
	if len(rr.nodes) == 0 {
		return nil
	}

	rr.mu.Lock()
	rr.currentIndex = (rr.currentIndex + 1) % len(rr.nodes)
	rr.mu.Unlock()

	return rr.nodes[rr.currentIndex]
}

func (rr *RoundRobin) CurrentIndex() int {
	return rr.currentIndex
}
