package node

import (
	"log"

	"github.com/baltej223/dukedb/internal/routing"
	"github.com/baltej223/dukedb/internal/storing"
	"github.com/baltej223/dukedb/internal/transport"
)

func handleGet(msg transport.ParsedMessage, me *Node) {
	sortedNodes := me.AllNodesSort()
	ownerNode := routing.FindOwner(msg.Key, sortedNodes)

	if ownerNode.NodeID == me.ID {
		if storing.Exists(msg.Key) {
			val, isOK := storing.Get(msg.Key)
			if !isOK {
				return
			}
			response := transport.CreateGetResponseMessage(msg.RequestID, true, val)

			peerToReply, ok := me.Cluster.GetPeer(msg.NodeID)
			if !ok {
				return
			}
			err := transport.SendMessage(peerToReply, response)
			if err != nil {
				return
			}
		}
	} else {
		// Here a new request needs to be made
	}
}

func handleGetResponse(msg transport.ParsedMessage, me *Node) {
	req, ok := me.GetPendingRequest(msg.RequestID)
	if !ok {
		log.Printf(
			"[node=%s] no pending request found for request_id=%s",
			me.ID,
			msg.RequestID,
		)
		return
	}
	req.ResultChan <- msg

	log.Printf(
		"[node=%s] pending request fulfilled request_id=%s",
		me.ID,
		msg.RequestID,
	)
	me.RemovePendingRequest(msg.RequestID)
}
