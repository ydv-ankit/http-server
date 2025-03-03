package main

import (
	"fmt"
	"net"
)

type ClientConnection struct {
	Conn     net.Conn
	Request  map[string]string
	Response map[string]string
}

func createConnection(c net.Conn) *ClientConnection {
	return &ClientConnection{
		Request:  make(map[string]string),
		Response: make(map[string]string),
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
	response := handleRequest(conn.Request["uri"])
	WriteTextResponse(&conn.Conn, "HTTP/1.1", 200, response)
}
