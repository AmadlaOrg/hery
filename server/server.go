package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const socketPath = "/tmp/hery.sock"

type Server interface {
	Start() error
}

// handleConnection
func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

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
	if _, err := os.Stat(socketPath); err == nil {
		err := os.Remove(socketPath)
		if err != nil {
			return
		}
	}

	l, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Println("Listen error:", err)
		return
	}
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {

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
