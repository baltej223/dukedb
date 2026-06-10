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

	log.Printf(
		"[ROUTE] op=PUT key=%s owner=%s ring=%v",
		msg.Key,
		keyOwner.NodeID,
		me.AllNodesSort(),
	)
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
}

func handlePutREJ(msg transport.ParsedMessage, me *Node) {
	req, ok := me.GetPendingRequest(msg.RequestID)
	if !ok {
		return
	}
	var finalResponse transport.ParsedMessage

	if msg.MembershipVersion > me.MembershipVersion {
		// schedule membership sync
		go SyncMembership(msg.NodeID, me, 20*time.Second)
	}
	if msg.SuggestedOwner == "" {
		finalResponse = msg
	} else {
		newPeerToTry := cluster.NewPeer(
			msg.SuggestedOwner,
			msg.SuggestedAddr,
		)
		if !me.Cluster.HasPeer(newPeerToTry.NodeID) {
			me.Cluster.AddPeer(newPeerToTry)
		}
		thisRequest, ok := me.GetPendingRequest(
			msg.RequestID,
		)
		if !ok {
			finalResponse = msg
		} else {
			retryResponse, err := me.SendRequestAndWait(
				newPeerToTry,
				thisRequest.Message,
				30*time.Second,
			)
			if err != nil {
				finalResponse = msg
			} else {
				finalResponse = retryResponse
			}
		}
	}

	req.ResultChan <- finalResponse
}
