package main

import (
	"fmt"
	"strings"
)

type Trie struct {
	Children     map[string]*Trie
	Handlers     map[string]func(*ClientConnection, map[string]string)
	IsDynamic    bool
	DynamicParam string
}

func NewTrie() *Trie {
	return &Trie{
		Children:     make(map[string]*Trie),
		Handlers:     make(map[string]func(*ClientConnection, map[string]string)),
		IsDynamic:    false,
		DynamicParam: "",
	}
}

var router = NewTrie()

func (t *Trie) AddRoute(method, path string, handler func(*ClientConnection, map[string]string)) {
	// check if path == "/"
	node := t
	if path == "/" {
		if _, exists := node.Children["/"]; exists {
			node.Handlers[method] = handler
			return
		}
		router.Children["/"] = NewTrie()
		node.Handlers[method] = handler
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
			node.DynamicParam = part[1:]
		} else {
			if _, exists := node.Children[part]; !exists {
				node.Children[part] = NewTrie()
			}
			node = node.Children[part]
		}
	}
	node.Handlers[method] = handler
}

func (t *Trie) FindRoute(path string, method string) (func(*ClientConnection, map[string]string), map[string]string, error) {
	node := t
	parts := strings.Split(path, "/")
	params := make(map[string]string)
	// if path == "/"
	if path == "/" {
		node = node.Children["/"]
	}
	for _, part := range parts {
		if part == "" {
			continue
		}
		if node.IsDynamic {
			params[node.DynamicParam] = part
			continue
		}
		fmt.Println("---------")
		fmt.Println(part)
		fmt.Println(node.Children[part])
		fmt.Println(node.Children[part].DynamicParam)
		fmt.Println(node.Children[part].Handlers)
		fmt.Println(node.Children[part].IsDynamic)
		fmt.Println("---------")
		_, exists := node.Children[part]
		if !exists {
			return nil, params, fmt.Errorf("route not found")
		}
		node = node.Children[part]
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
		fmt.Println("params: ", t.Children[key].DynamicParam)
		fmt.Println("handlers: ", t.Children[key].Handlers)
		t.Children[key].getRoutes()
	}
}

func (c *ClientConnection) handleRequest() {
	router.AddRoute("GET", "/files/:filename", serveStaticFiles)
	router.AddRoute("GET", "/hello", sayHello)
	fmt.Println("Requested Path:", c.Request.PATH)

	handler, params, err := router.FindRoute(c.Request.PATH, c.Request.METHOD)
	if err != nil {
		fmt.Println(err)
		return
	}
	if handler != nil {
		fmt.Println("params", params)
		fmt.Println("calling handler")
		handler(c, params)
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
