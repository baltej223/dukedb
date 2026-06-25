package node

import (
	"github.com/baltej223/dukedb/internal/transport"
	dukelog "github.com/baltej223/dukedb/log"
)

func handlePong(msg transport.ParsedMessage, me *Node) {
	dukelog.Printf(
		"[node=%s] received PONG request_id=%s from node=%s",
		me.ID,
		msg.RequestID,
		msg.NodeID,
	)

	req, ok := me.GetPendingRequest(msg.RequestID)
	if !ok {
		dukelog.Printf(
			"[node=%s] no pending request found for request_id=%s",
			me.ID,
			msg.RequestID,
		)
		return
	}

	dukelog.Printf(
		"[node=%s] fulfilling pending request request_id=%s",
		me.ID,
		msg.RequestID,
	)

	req.ResultChan <- msg

	dukelog.Printf(
		"[node=%s] pending request fulfilled request_id=%s",
		me.ID,
		msg.RequestID,
	)
}
