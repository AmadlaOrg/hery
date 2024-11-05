package client

import (
	"bufio"
	"fmt"
	"github.com/AmadlaOrg/hery/server"
	"log"
	"net"
	"os"
)

// IClient
type IClient interface {
	Connect()
}

// SClient
type SClient struct{}

// Connect
func (s *SClient) Connect() {
	conn, err := net.Dial("unix", server.SocketPath)
	if err != nil {
		log.Fatalf("Dial error: %v", err)
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("Close error:", err)
		}
	}(conn)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter commands (type 'exit' to quit):")

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		if text == "exit\n" {
			break
		}
		_, err := conn.Write([]byte(text))
		if err != nil {
			log.Printf("Write error: %v", err)
			break
		}
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		log.Printf("Response: %v", response)
	}
}
