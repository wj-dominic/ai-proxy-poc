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
	config, err := config.LoadConfig("")
	if err != nil {
		log.Fatal(err)
	}

	nodes := make([]*node.Node, 0)
	for _, n := range config.Nodes {
		parsedURL, err := url.Parse(n.URL)
		if err != nil {
			log.Fatal(err)
		}

		nodes = append(nodes, node.NewNode(n.ID, parsedURL, n.MaxBPM, n.MaxRPM))
	}

	server := &http.Server{
		Handler: loadbalancer.NewLoadBalancer(nodes, loadbalancer.NewRoundRobin(nodes)),
		Addr:    ":8081",
	}

	go func() {
		log.Println("Server is running on port", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
