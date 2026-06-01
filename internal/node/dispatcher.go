package node

import (
	"log"

	"github.com/baltej223/dukedb/internal/transport"
)

func Dispatch(msg transport.ParsedMessage, me *Node) {
	log.Printf(
		"[node=%s] dispatching %s request_id=%s",
		me.ID,
		msg.Type,
		msg.RequestID,
	)

	switch msg.Type {
	case transport.PING:
		handlePing(msg, me)
	case transport.PONG:
		handlePong(msg, me)
		// case transport.JOIN:
		// 	handlers.handleJoin(msg)
	}
}
