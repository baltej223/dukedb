// Package transport: Defines transport layer.
package transport

import (
	"bufio"
	"fmt"
	"net"
)

type Message struct {
	Type    string
	From    string
	Payload []byte
}

type Server struct {
	address string
}

func NewServer(address string) *Server {
	return &Server{
		address: address,
	}
}

func (s *Server) Start(connectionHandler func(conn net.Conn)) error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	defer listener.Close()

	fmt.Println("tcp server listening on", s.address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		// go HandleConnection(conn)
		go connectionHandler(conn)
	}
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()

	fmt.Println("new connection from", remoteAddr)

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("connection closed:", remoteAddr)
			return
		}

		fmt.Printf("received from %s: %s", remoteAddr, message)

		response := "ack: " + message

		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("write error:", err)
			return
		}
	}
}
