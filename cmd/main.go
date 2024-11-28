package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/wj-dominic/ai-proxy-poc/pkg/config"
	loadbalancer "github.com/wj-dominic/ai-proxy-poc/pkg/load-balancer"
	"github.com/wj-dominic/ai-proxy-poc/pkg/node"
)

func main() {
	// 1. Load configuration
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// 2. Create nodes
	nodes := make([]*node.Node, 0)
	for _, n := range config.Nodes {
		parsedURL, err := url.Parse(n.URL)
		if err != nil {
			log.Fatal(err)
		}

		nodes = append(nodes, node.NewNode(n.ID, parsedURL, n.MaxRPM, n.MaxBPM))
	}

	// 3. Choose load balancing algorithm
	// TODO: Implement other algorithms
	var balancingAlgorithm loadbalancer.Algorithm
	switch config.LB.Algorithm {
	case "round-robin":
		log.Println("Using Round Robin algorithm")
		balancingAlgorithm = loadbalancer.NewRoundRobin(nodes)
	default:
		log.Println("Using Round Robin algorithm as default")
		balancingAlgorithm = loadbalancer.NewRoundRobin(nodes)
	}

	// 4. Create server with lb
	server := &http.Server{
		Handler: loadbalancer.NewLoadBalancer(nodes, balancingAlgorithm),
		Addr:    ":8081",
	}

	// 5. Run server
	go func() {
		log.Println("Server is running on port", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// 6. Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
