package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strconv"
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

func readFileContent(path string) ([]byte, string, error) {
	// Check if file exists
	fileInfo, err := os.Stat(path)
	if err != nil || fileInfo.IsDir() {
		return nil, "", fmt.Errorf("file not found")
	}

	// Read file content
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("unable to read file")
	}

	// Get MIME type dynamically
	ext := filepath.Ext(path) // Extract file extension (e.g., .html, .jpg)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream" // Fallback for unknown types
	}

	return content, mimeType, nil
}

func serveStaticFiles(c *ClientConnection, params map[string]string) {
	// Construct the full file path
	filePath := "./static/" + params["filename"]

	fmt.Println("got filepath", filePath)
	// Read file content dynamically
	content, mimeType, err := readFileContent(filePath)
	if err != nil {
		fmt.Println(err)
		// File not found → Return 404
		c.Response = Response{
			STATUS:   404,
			PROTOCOL: "HTTP/1.1",
			HEADERS: map[string]string{
				"Content-Type":   "text/plain",
				"Content-Length": "0",
			},
		}
		c.WriteTextResponse()
		return
	}

	// File found → Return file with correct headers
	c.Response = Response{
		STATUS:   200,
		PROTOCOL: "HTTP/1.1",
		HEADERS: map[string]string{
			"Content-Type":   mimeType,
			"Content-Length": strconv.Itoa(len(content)),
		},
		BODY: content,
	}
	// Compress content with gzip if client supports it
	// brotli
	if strings.Contains(c.Request.HEADERS["Accept-Encoding"], "gzip") && strings.Contains(mimeType, "text/") {
		var buffer bytes.Buffer
		zw := gzip.NewWriter(&buffer)
		_, err := zw.Write(content)
		if err != nil {
			fmt.Println("Error compressing content:", err)
			return
		}
		err = zw.Close()
		if err != nil {
			fmt.Println("Error closing gzip writer:", err)
			return
		}
		c.Response.HEADERS["Content-Encoding"] = "gzip"
		c.Response.BODY = buffer.Bytes()
		c.Response.HEADERS["Content-Length"] = strconv.Itoa(len(c.Response.BODY))
	}
	c.WriteTextResponse()
}

func sayHello(c *ClientConnection, params map[string]string) {
	c.Response = Response{
		PROTOCOL: "http/1.1",
		STATUS:   200,
		HEADERS: map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": "5",
		},
		BODY: []byte("Hello"),
	}
	c.WriteTextResponse()
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
