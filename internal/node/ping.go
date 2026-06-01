package node

import (
	"log"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/transport"
)

// func handlePing(msg transport.ParsedMessage, me Node) {
// 	pong := transport.CreatePongMessage(
// 		msg.RequestID,
// 		me.ID,
// 	)
//
// 	peer, _ := cluster.PeerFromNodeID(msg.NodeID)
//
// 	transport.SendMessage(peer, pong)
// }

func handlePing(msg transport.ParsedMessage, me Node) {
	log.Printf(
		"[node=%s] received PING request_id=%s from node=%s",
		me.ID,
		msg.RequestID,
		msg.NodeID,
	)

	pong := transport.CreatePongMessage(
		msg.RequestID,
		me.ID,
	)

	peer, err := cluster.PeerFromNodeID(msg.NodeID)
	if err != nil {
		log.Printf(
			"[node=%s] failed to find peer %s: %v",
			me.ID,
			msg.NodeID,
			err,
		)
		return
	}

	log.Printf(
		"[node=%s] sending PONG request_id=%s to node=%s",
		me.ID,
		msg.RequestID,
		peer.NodeID,
	)

	if err := transport.SendMessage(peer, pong); err != nil {
		log.Printf(
			"[node=%s] failed to send PONG to node=%s: %v",
			me.ID,
			peer.NodeID,
			err,
		)
		return
	}

	log.Printf(
		"[node=%s] PONG sent successfully request_id=%s",
		me.ID,
		msg.RequestID,
	)
}
