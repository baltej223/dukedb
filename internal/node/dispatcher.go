package node

import (
	"github.com/baltej223/dukedb/internal/transport"
)

func Dispatch(msg transport.ParsedMessage, me Node) {
	switch msg.Type {
	case transport.PING:
		handlePing(msg, me)
		//
		// case transport.PONG:
		// 	handlers.handlePong(msg)
		//
		// case transport.JOIN:
		// 	handlers.handleJoin(msg)
	}
}
