package client

import (
	"bufio"
	"fmt"
	"github.com/AmadlaOrg/hery/server"
	"net"
	"os"
)

type Client interface {
	Connect() error
}

// Connect
func Connect() {
	conn, err := net.Dial("unix", server.SocketPath)
	if err != nil {
		fmt.Println("Dial error:", err)
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

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
			fmt.Println("Write error:", err)
			break
		}
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		fmt.Println("Response:", response)
	}
}
