package loadbalancer

import (
	"net/http"

	"github.com/wj-dominic/ai-proxy-poc/pkg/node"
)

type Algorithm interface {
	NextNode() *node.Node
}

type LoadBalancer struct {
	nodes     []*node.Node
	algorithm Algorithm
}

func NewLoadBalancer(nodes []*node.Node, algorithm Algorithm) *LoadBalancer {
	if nodes == nil {
		nodes = make([]*node.Node, 0)
	}

	if algorithm == nil {
		algorithm = NewRoundRobin(nodes)
	}

	return &LoadBalancer{
		nodes:     nodes,
		algorithm: algorithm,
	}
}

func (lb *LoadBalancer) NextNode() *node.Node {
	return lb.algorithm.NextNode()
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	node := lb.NextNode()
	if node == nil {
		http.Error(w, "No nodes available", http.StatusServiceUnavailable)
		return
	}

	if !node.IsAllowRequest(r.ContentLength) {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	// 서버에 문제가 발생한 경우는 노드의 rate limit을 어떻게 처리하지?
	node.ServeHTTP(w, r)
}
