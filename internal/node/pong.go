package node

import (
	"log"

	"github.com/baltej223/dukedb/internal/transport"
)

func handlePong(msg transport.ParsedMessage, me *Node) {
	log.Printf(
		"[node=%s] received PONG request_id=%s from node=%s",
		me.ID,
		msg.RequestID,
		msg.NodeID,
	)

	req, ok := me.GetPendingRequest(msg.RequestID)
	if !ok {
		log.Printf(
			"[node=%s] no pending request found for request_id=%s",
			me.ID,
			msg.RequestID,
		)
		return
	}

	log.Printf(
		"[node=%s] fulfilling pending request request_id=%s",
		me.ID,
		msg.RequestID,
	)

	req.ResultChan <- msg

	log.Printf(
		"[node=%s] pending request fulfilled request_id=%s",
		me.ID,
		msg.RequestID,
	)
	me.RemovePendingRequest(msg.RequestID)
}
