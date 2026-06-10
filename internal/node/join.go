package node

import (
	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/transport"
)

func handleJoin(msg transport.ParsedMessage, me *Node) {
	newPeer := cluster.NewPeer(
		msg.NodeID,
		msg.Addr,
	)
	me.Cluster.AddPeer(newPeer)
	me.MembershipVersion++
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
		return
	}
	if msg.MembershipVersion > me.MembershipVersion {
		me.MembershipVersion = msg.MembershipVersion
	}

	req.ResultChan <- msg

	for _, value := range msg.Peers {
		if value.NodeID != me.ID {
			me.Cluster.AddPeer(value)
		}
	}
}
