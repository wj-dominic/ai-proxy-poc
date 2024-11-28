package node

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Node struct {
	id  string
	url *url.URL

	maxRPM     int
	currentRPM int

	maxBPM     int64
	currentBPM int64

	proxy *httputil.ReverseProxy

	mu sync.Mutex
}

func NewNode(id string, url *url.URL, maxRPM int, maxBPM int64) *Node {
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Director = dialector(url)

	return &Node{
		id:         id,
		url:        url,
		maxRPM:     maxRPM,
		currentRPM: 0,
		maxBPM:     maxBPM,
		currentBPM: 0,
		proxy:      proxy,
	}
}

func (n *Node) ID() string {
	return n.id
}

func (n *Node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n.proxy.ServeHTTP(w, r)
}

func (n *Node) IsAllowRequest(bodySize int64) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Check RPM, BPM
	if n.currentRPM+1 > n.maxRPM || n.currentBPM+bodySize > n.maxBPM {
		return false
	}

	n.currentRPM++
	n.currentBPM += bodySize

	return true
}

func dialector(url *url.URL) func(*http.Request) {
	return func(req *http.Request) {
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Origin-Host", url.Host)
		req.Host = url.Host
		req.URL.Host = url.Host
		req.URL.Scheme = url.Scheme
	}
}
