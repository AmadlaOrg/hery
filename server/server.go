package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Server interface {
	Start() error
}

// handleConnection
func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}(conn)
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		response := processCommand(message)
		_, err = conn.Write([]byte(response + "\n"))
		if err != nil {
			return
		}
	}
}

// processCommand
func processCommand(cmd string) string {
	// Here you can add logic to process the command
	return "Processed: " + cmd
}

// Start
func Start() {
	if _, err := os.Stat(SocketPath); err == nil {
		err := os.Remove(SocketPath)
		if err != nil {
			return
		}
	}

	l, err := net.Listen("unix", SocketPath)
	if err != nil {
		fmt.Println("Listen error:", err)
		return
	}
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Fatalf("Error closing socket: %v", err)
		}
	}(l)

	// Handle SIGINT and SIGTERM for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		err := l.Close()
		if err != nil {
			return
		}
		os.Exit(0)
	}()

	fmt.Println("Daemon started. Waiting for connections...")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go handleConnection(conn)
	}
}
