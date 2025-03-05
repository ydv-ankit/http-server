package app

import (
	"fmt"
	"strings"
)

type Trie struct {
	Children      map[string]*Trie
	Handlers      map[string][]func(*ClientConnection, map[string]string)
	IsDynamic     bool
	DynamicParams []string
}

func NewTrie() *Trie {
	return &Trie{
		Children:      make(map[string]*Trie),
		Handlers:      make(map[string][]func(*ClientConnection, map[string]string)),
		IsDynamic:     false,
		DynamicParams: nil,
	}
}

var router = NewTrie()

func (t *Trie) AddRoute(method, path string, handlers ...func(*ClientConnection, map[string]string)) {
	// check if path == "/"
	node := t
	if path == "/" {
		node.Handlers[method] = append(node.Handlers[method], handlers...)
		return
	}
	// Split path into parts
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if part == "" {
			continue
		}
		if part[0] == ':' {
			node.IsDynamic = true
			node.DynamicParams = append(node.DynamicParams, part[1:])
		} else {
			if _, exists := node.Children[part]; !exists {
				node.Children[part] = NewTrie()
			}
			node = node.Children[part]
		}
	}
	node.Handlers[method] = append(node.Handlers[method], handlers...)
}

func (t *Trie) FindRoute(path string, method string) ([]func(*ClientConnection, map[string]string), map[string]string, error) {
	node := t
	params := make(map[string]string)
	parts := strings.Split(path, "/")
	paramIdx := 0
	for _, part := range parts {
		if part == "" {
			continue
		}
		fmt.Println(part)
		fmt.Println("node", node)
		if node.IsDynamic {
			if len(node.DynamicParams) > paramIdx {
				params[node.DynamicParams[paramIdx]] = part
				paramIdx++
			}
		}
		_, exists := node.Children[part]
		if exists {
			node = node.Children[part]
		}
	}
	if node.Handlers[method] != nil {
		return node.Handlers[method], params, nil
	}
	return nil, params, fmt.Errorf("route not found")
}

func (t *Trie) getRoutes() {
	if len(t.Children) == 0 {
		return
	}
	for key := range t.Children {
		fmt.Println("route: ", key)
		fmt.Println("isDynamic: ", t.Children[key].IsDynamic)
		fmt.Println("params: ", t.Children[key].DynamicParams)
		fmt.Println("handlers: ", t.Children[key].Handlers)
		fmt.Println("children: ", t.Children[key].Children)
		t.Children[key].getRoutes()
	}
}

func (c *ClientConnection) handleRequest() {
	router.AddRoute("GET", "/files/:filename", serveStaticFiles)
	router.AddRoute("GET", "/user/:name/:filename", checkAuth, serveStaticFiles)
	fmt.Println("Requested Path:", c.Request.PATH)

	router.getRoutes()
	handlers, params, err := router.FindRoute(c.Request.PATH, c.Request.METHOD)
	if err != nil {
		fmt.Println(err)
		return
	}
	if handlers != nil || len(handlers) == 0 {
		fmt.Println("params", params)
		fmt.Println("calling handlers", handlers)
		for _, handler := range handlers {
			handler(c, params)
		}
		return
	}
	fmt.Println("no handler attached")
	c.Response = Response{
		STATUS:   400,
		PROTOCOL: "HTTP/1.1",
		HEADERS: map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": "0",
		},
	}
	c.WriteTextResponse()
}
