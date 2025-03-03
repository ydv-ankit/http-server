package main

import (
	"fmt"
	"net"
)

func startServer() {
	conn, err := net.Listen("tcp", "0.0.0.0:4444")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fmt.Println("Server is listening...")
	// listen to client calls infinitely
	for {
		c, err := conn.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}
