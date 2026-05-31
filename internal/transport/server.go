// Package transport: Defines transport layer.
package transport

import (
	"fmt"
	"net"

	"github.com/baltej223/dukedb/internal/cluster"
)

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

func HandleMessage(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()

	fmt.Println("new connection from", remoteAddr)

	for {
		message, err := readMessage(conn)
		if err != nil {
			fmt.Println("connection closed:", remoteAddr)
			return
		}

		fmt.Printf("received from %s: %s", remoteAddr, message)

		messageParsed, err := Parse(message)

		fmt.Println(messageParsed.String())
		if err != nil {
			fmt.Println("Error in message reading, Err: %s", err)
			continue
		}

		if messageParsed.Type == PING {
			pong := CreatePongMessage(messageParsed.RequestID, "a")
			peer, _ := cluster.PeerFromNodeID(messageParsed.NodeID)
			_ = SendMessage(peer, pong, "SENDING:")

		}
	}
}
