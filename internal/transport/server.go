// Package transport: Defines transport layer.
package transport

import (
	"net"

	dukelog "github.com/baltej223/dukedb/log"
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

	dukelog.Println("tcp server listening on", s.address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			dukelog.Println("accept error:", err)
			continue
		}

		// go HandleConnection(conn)
		go connectionHandler(conn)
	}
}

func HandleConnection(conn net.Conn, dispatch func(ParsedMessage)) {
	dukelog.Println("New connection")

	raw, err := readMessage(conn)
	if err != nil {
		dukelog.Println("read error:", err)
		return
	}

	// dukelog.Printf("raw message:\n%s", raw)

	parsed, err := Parse(raw)
	if err != nil {
		dukelog.Println("parse error:", err)
		return
	}

	dukelog.Printf(
		"Parsed Type=%s Request_id=%s",
		parsed.Type,
		parsed.RequestID,
	)

	dispatch(parsed)
}
