// Package transport: Defines transport layer.
package transport

import (
	"fmt"
	"net"
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

func HandleConnection(conn net.Conn, dispatch func(ParsedMessage)) {
	raw, err := readMessage(conn)
	if err != nil {
		return
	}

	parsed, err := Parse(raw)
	if err != nil {
		return
	}

	dispatch(parsed)
}
