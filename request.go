package main

import (
	"fmt"
	"strings"
)

func splitRequest(data []byte) (string, string, string) {
	request := string(data)
	if len(strings.Split(request, "\r\n")) > 1 {
		startLine := strings.Split(request, "\r\n")[0]
		if len(strings.Split(request[len(startLine):], "\r\n\r\n")) == 1 {
			headers := strings.Split(request[len(startLine):], "\r\n\r\n")[0]
			return startLine, headers, ""
		}
		if len(strings.Split(request[len(startLine):], "\r\n\r\n")) >= 2 {
			headers := strings.Split(request[len(startLine):], "\r\n\r\n")[0]
			body := strings.Split(request[len(startLine):], "\r\n\r\n")[1]
			return startLine, headers, body
		}
		return startLine, "", ""
	} else {
		return "", "", ""
	}
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

func (c *ClientConnection) parseStartLine(startLine string) {
	parts := strings.Split(startLine, " ")
	if len(parts) == 3 {
		c.Request.METHOD = parts[0]
		c.Request.PATH = parts[1]
		c.Request.VERSION = parts[2]
	}
}

func (c *ClientConnection) parseRequest(data []byte) error {
	startLine, headers, body := splitRequest(data)
	if startLine == "" {
		return fmt.Errorf("empty request")
	}
	c.parseStartLine(startLine)
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
