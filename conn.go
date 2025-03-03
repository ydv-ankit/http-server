package main

import (
	"fmt"
	"net"
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
	BODY     string
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
	data := make([]byte, 1024)
	n, err := conn.Conn.Read(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = conn.parseRequest(data[:n])
	if err != nil {
		WriteTextResponse(&conn.Conn, "HTTP/1.1", 400, "")
		return
	}
	conn.showRequest()
	response := handleRequest(conn.Request.URI)
	WriteTextResponse(&conn.Conn, "HTTP/1.1", 200, response)
}
