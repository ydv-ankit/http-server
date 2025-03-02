package main

import (
	"fmt"
	"net"
	"strings"
)

func handleConnection(c net.Conn) {
	defer c.Close()
	fmt.Println("Accepted connection from", c.RemoteAddr())
	data := make([]byte, 4096)
	n, _ := c.Read(data)
	fmt.Println(data[:n])
	firstLine := strings.Split(strings.Split(string(data[:n]), "\r\n")[0], " ")
	reqMethod := firstLine[0]
	reqURI := firstLine[1]
	reqProto := firstLine[2]
	fmt.Println("Method:", reqMethod)
	fmt.Println("URI:", reqURI)
	fmt.Println("Protocol:", reqProto)
	WriteResponse(&c, 200, "hello there")
}
