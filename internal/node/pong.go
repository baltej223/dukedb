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

	// if me.IsSuspectedDead(msg.NodeID) {
	// 	me.RemoveSuspectedDeadPeer(msg.NodeID)
	// }

	dukelog.Printf(
		"[node=%s] pending request fulfilled request_id=%s",
		me.ID,
		msg.RequestID,
	)
}
