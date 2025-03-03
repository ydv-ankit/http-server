package main

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strconv"
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

func (c *ClientConnection) handleRequest() {
	fmt.Println("Requested Path:", c.Request.PATH)

	// Construct the full file path
	filePath := "./static" + c.Request.PATH

	// Read file content dynamically
	content, mimeType, err := readFileContent(filePath)
	if err != nil {
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
	c.WriteTextResponse()
}
