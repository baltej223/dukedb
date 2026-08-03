package node

import (
	"github.com/baltej223/dukedb/internal/transport"
	dukelog "github.com/baltej223/dukedb/log"
)

func handlePing(msg transport.ParsedMessage, me *Node) {
	pong := transport.CreatePongMessage(
		msg.RequestID,
		me.ID,
	)

	peer, ok := me.Cluster.GetPeer(
		msg.NodeID,
	)
	if !ok {
		dukelog.Printf(
			"[node=%s] failed to find peer %s",
			me.ID,
			msg.NodeID,
		)
	}

	if err := transport.SendMessage(peer, pong); err != nil {
		dukelog.Printf(
			"[node=%s] failed to send PONG to node=%s: %v",
			me.ID,
			peer.NodeID,
			err,
		)
		return
	}

	dukelog.Printf(
		"[node=%s] PONG sent successfully request_id=%s",
		me.ID,
		msg.RequestID,
	)
}
