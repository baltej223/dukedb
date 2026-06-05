package node

import (
	"log"
	"time"

	"github.com/baltej223/dukedb/internal/cluster"
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

func handleGetREJ(msg transport.ParsedMessage, me *Node) {
	req, ok := me.GetPendingRequest(msg.RequestID)
	if !ok {
		return
	}
	var finalResponse transport.ParsedMessage

	if msg.MembershipVersion > me.MembershipVersion {
		// schedule membership sync
	}
	if msg.SuggestedOwner != "" {
		newPeerToTry := cluster.NewPeer(msg.SuggestedOwner, msg.SuggestedAddr)
		if !(me.Cluster.HasPeer(newPeerToTry.NodeID)) {
			me.Cluster.AddPeer(newPeerToTry)
		}

		thisRequest, ok := me.GetPendingRequest(msg.RequestID)
		if !ok {
			return
		}
		newMessage, err := transport.CreateGetMessage(
			thisRequest.Message.Headers["KEY"],
			me.ID,
			me.MembershipVersion,
		)
		if err != nil {
			finalResponse = msg
		}

		getResponseFromNewPeer, err := me.SendRequestAndWait(newPeerToTry, newMessage, 30*time.Second)
		if err != nil {
			finalResponse = msg
		}
		finalResponse = getResponseFromNewPeer
	} else {
		finalResponse = msg
	}

	req.ResultChan <- finalResponse

	me.RemovePendingRequest(msg.RequestID)
}
