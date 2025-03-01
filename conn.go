package main

import (
	"fmt"
	"net"
)

func handleConnection(c net.Conn) {
	defer c.Close()
	fmt.Println("Accepted connection from", c.RemoteAddr())
}
