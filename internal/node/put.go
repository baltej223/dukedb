package node

import (
	"log"
	"time"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/routing"
	"github.com/baltej223/dukedb/internal/storing"
	"github.com/baltej223/dukedb/internal/transport"
)

func handlePut(msg transport.ParsedMessage, me *Node) {
	log.Printf(
		"[node=%s] PUT ENTER request_id=%s key=%s sender=%s",
		me.ID,
		msg.RequestID,
		msg.Key,
		msg.NodeID,
	)

	key := msg.Key

	keyOwner := routing.FindOwner(key, me.AllNodesSort())
	if keyOwner.NodeID == me.ID {

		storing.Put(key, msg.Value)

		response := transport.CreatePutACKMessage(msg.RequestID)

		peerToReply, ok := me.Cluster.GetPeer(msg.NodeID)
		if !ok {
			return
		}
		err := transport.SendMessage(peerToReply, response)
		if err != nil {
			return
		}

	} else {
		response := transport.CreatePutREJMessage(
			msg.RequestID,
			keyOwner,
			me.MembershipVersion,
		)

		peerToReply, ok := me.Cluster.GetPeer(msg.NodeID)
		if !ok {
			return
		}

		err := transport.SendMessage(peerToReply, response)
		if err != nil {
			return
		}

	}
}

func handlePutACK(msg transport.ParsedMessage, me *Node) {
	req, ok := me.GetPendingRequest(msg.RequestID)
	if !ok {
		return
	}
	req.ResultChan <- msg
	me.RemovePendingRequest(msg.RequestID)
}

func handlePutREJ(msg transport.ParsedMessage, me *Node) {
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
			finalResponse = msg
		}
		newMessage := thisRequest.Message

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
