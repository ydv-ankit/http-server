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
