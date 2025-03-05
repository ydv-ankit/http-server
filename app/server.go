package app

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func StartServer() {
	conn, err := net.Listen("tcp", "0.0.0.0:4444")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down server gracefully...")
		conn.Close()
		os.Exit(0)
	}()

	// listen to client calls infinitely
	fmt.Println("Server is listening...")
	for {
		c, err := conn.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}
