package transport

import (
	"net"

	"github.com/baltej223/dukedb/internal/node"
)

// Client, the function which actually communicates with the other nodes.

func Send(n *node.Node, s string) (string, error) {
	conn, err := net.Dial("tcp", n.GetFullHostname())
	if err != nil {
		return "", err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(s + "\n"))
	if err != nil {
		return "", err
	}

	buffer := make([]byte, 1024)

	nRead, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}

	response := string(buffer[:nRead])

	return response, nil
}
