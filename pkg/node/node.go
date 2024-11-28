package node

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Node struct {
	id     string
	url    *url.URL
	maxBPM int
	maxRPM int

	proxy *httputil.ReverseProxy
}

func NewNode(id string, url *url.URL, maxBPM int, maxRPM int) *Node {
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Director = dialector(url)

	return &Node{
		id:     id,
		url:    url,
		maxBPM: maxBPM,
		maxRPM: maxRPM,
		proxy:  proxy,
	}
}

func (n *Node) ID() string {
	return n.id
}

func (n *Node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n.proxy.ServeHTTP(w, r)
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
