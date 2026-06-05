package node

import (
	"log"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/transport"
)

func handleJoin(msg transport.ParsedMessage, me *Node) {
	newPeer := cluster.NewPeer(
		msg.NodeID,
		msg.Addr,
	)
	me.Cluster.AddPeer(newPeer)

	// Send the JOIN_ACK message here
	currentPeer := me.Cluster.GetPeers()
	joinACK := transport.CreateJoinACKMessage(
		msg.RequestID,
		currentPeer,
		me.MembershipVersion,
	)
	transport.SendMessage(newPeer, joinACK)
}

func handleJoinACK(msg transport.ParsedMessage, me *Node) {
	req, ok := me.GetPendingRequest(msg.RequestID)
	if !ok {
		log.Printf(
			"[node=%s] no pending request found for request_id=%s",
			me.ID,
			msg.RequestID,
		)
		return
	}
	if msg.MembershipVersion > me.MembershipVersion {
		me.MembershipVersion = msg.MembershipVersion
	}

	req.ResultChan <- msg

	for _, value := range msg.Peers {
		me.Cluster.AddPeer(value)
	}

	log.Printf(
		"[node=%s] pending request fulfilled request_id=%s",
		me.ID,
		msg.RequestID,
	)
	me.RemovePendingRequest(msg.RequestID)
}
