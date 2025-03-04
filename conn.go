package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	METHOD  string
	URI     string
	VERSION string
	HEADERS map[string]string
	PATH    string
	BODY    string
}

type Response struct {
	PROTOCOL string
	STATUS   int
	HEADERS  map[string]string
	BODY     []byte
}

type ClientConnection struct {
	Conn     net.Conn
	Request  Request
	Response Response
}

func createConnection(c net.Conn) *ClientConnection {
	return &ClientConnection{
		Request:  Request{},
		Response: Response{},
		Conn:     c,
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	conn := createConnection(c)

	for {
		data := make([]byte, 1024)
		n, err := conn.Conn.Read(data)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = conn.parseRequest(data[:n])
		if err != nil {
			fmt.Println(err)
			return
		}

		conn.handleRequest()

		// Default timeout (if Keep-Alive timeout is missing or invalid)
		timeout := 5

		// Check for Keep-Alive header and extract timeout if present
		if conn.Request.HEADERS["Connection"] == "keep-alive" {
			if keepAliveVal, exists := conn.Request.HEADERS["Keep-Alive"]; exists {
				parts := strings.Split(keepAliveVal, ",")
				for _, part := range parts {
					part = strings.TrimSpace(part)
					if strings.HasPrefix(part, "timeout=") {
						val := strings.TrimPrefix(part, "timeout=")
						if parsedTimeout, err := strconv.Atoi(val); err == nil {
							timeout = parsedTimeout
						}
						break
					}
				}
			}
		} else {
			break // No Keep-Alive, close the connection
		}

		// Set connection deadline based on extracted timeout
		conn.Conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}
}
