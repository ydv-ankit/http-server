package main

import (
	"fmt"
	"strings"
)

func (c *ClientConnection) parseRequest(data []byte) error {
	request := string(data)
	startLine := strings.Split(request, "\r\n")[0]
	fmt.Println(len(startLine))
	headers := strings.Split(request[len(startLine):], "\r\n\r\n")
	fmt.Println("headers[0]", headers[0])

	return fmt.Errorf("bad request")
}
