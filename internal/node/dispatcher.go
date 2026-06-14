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
		go handlePing(msg, me)
	case transport.PONG:
		go handlePong(msg, me)
	case transport.JOIN:
		go handleJoin(msg, me)
	case transport.JOIN_ACK:
		go handleJoinACK(msg, me)
	case transport.GOSSIPMEMBERSHIP:
		go handleMembership(msg, me)
	case transport.GET:
		go handleGet(msg, me)
	case transport.GET_RESPONSE:
		go handleGetResponse(msg, me)
	case transport.GET_REJ:
		go handleGetREJ(msg, me)
	case transport.PUT:
		go handlePut(msg, me)
	case transport.PUT_ACK:
		go handlePutACK(msg, me)
	case transport.PUT_REJ:
		go handlePutREJ(msg, me)
	case transport.SYNC_MEMBERSHIP:
		go handleSYNCMembership(msg, me)
	case transport.SYNC_MEMBERSHIP_RESPONSE:
		go handleSYNCMembershipResponse(msg, me)
	}
}
