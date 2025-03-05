package app

import (
	"fmt"
	"strings"
)

func splitRequest(data []byte) (string, string, string) {
	request := string(data)

	// Find the double CRLF that separates headers and body
	separatorIndex := strings.Index(request, "\r\n\r\n")

	// Split request into start line, headers, and body
	lines := strings.Split(request[:separatorIndex], "\r\n")
	if len(lines) == 0 {
		return "", "", ""
	}

	startLine := lines[0]
	headers := strings.Join(lines[1:], "\r\n")

	// If no headers or body, return empty
	if separatorIndex == -1 {
		return startLine, headers, ""
	}
	body := request[separatorIndex+4:] // Skip "\r\n\r\n"

	return startLine, headers, body
}

func (c *ClientConnection) parseHeaders(headers string) {
	c.Request.HEADERS = make(map[string]string)
	for _, header := range strings.Split(headers, "\r\n") {
		parts := strings.Split(header, ": ")
		if len(parts) == 2 {
			c.Request.HEADERS[parts[0]] = parts[1]
		}
	}
}

func (c *ClientConnection) parseStartLine(startLine string) error {
	parts := strings.Split(startLine, " ")
	if len(parts) == 3 {
		c.Request.METHOD = parts[0]
		c.Request.PATH = parts[1]
		c.Request.VERSION = parts[2]
	} else {
		fmt.Println("Invalid start line")
		return fmt.Errorf("invalid start line")
	}
	return nil
}

func (c *ClientConnection) parseRequest(data []byte) error {
	startLine, headers, body := splitRequest(data)
	if startLine == "" {
		return fmt.Errorf("empty request")
	}
	err := c.parseStartLine(startLine)
	if err != nil {
		return err
	}
	c.parseHeaders(headers)
	c.Request.BODY = body
	return nil
}

func (c *ClientConnection) showRequest() {
	fmt.Println("Method: ", c.Request.METHOD)
	fmt.Println("Path: ", c.Request.PATH)
	fmt.Println("Version: ", c.Request.VERSION)
	fmt.Println("Headers: ")
	for key, value := range c.Request.HEADERS {
		fmt.Println(key, ": ", value)
	}
	fmt.Println("Body: ", c.Request.BODY)
}
