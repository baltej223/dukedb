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
	case transport.JOIN:
		handleJoin(msg, me)
	case transport.JOIN_ACK:
		handleJoinACK(msg, me)
	case transport.GOSSIPMEMBERSHIP:
		handleMembership(msg, me)
	case transport.GET:
		handleGet(msg, me)
	case transport.GET_RESPONSE:
		handleGetResponse(msg, me)
	case transport.GET_REJ:
		handleGetREJ(msg, me)
	case transport.PUT:
		handlePut(msg, me)
	case transport.PUT_ACK:
		handlePutACK(msg, me)
	case transport.PUT_REJ:
		handlePutREJ(msg, me)
	case transport.SYNC_MEMBERSHIP:
		handleSYNCMembership(msg, me)
	case transport.SYNC_MEMBERSHIP_RESPONSE:
		handleSYNCMembershipResponse(msg, me)
	}
}
