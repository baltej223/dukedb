package node

import (
	"time"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/routing"
	"github.com/baltej223/dukedb/internal/storing"
	"github.com/baltej223/dukedb/internal/transport"
)

func handleGet(
	msg transport.ParsedMessage,
	me *Node,
) {
	sortedNodes := me.AllNodesSort()

	ownerNode := routing.FindOwner(
		msg.Key,
		sortedNodes,
	)
	reject := func(requestID string, me *Node) {
		response := transport.CreateGetREJMessage(requestID, ownerNode, me.MembershipVersion)
		peerToReply, ok := me.Cluster.GetPeer(
			msg.NodeID,
		)
		if !ok {
			return
		}

		err := transport.SendMessage(
			peerToReply,
			response,
		)
		if err != nil {
			return
		}
	}
	if (ownerNode.NodeID == me.ID) || IsReplica(msg.Key, me) {
		if storing.Exists(msg.Key) {

			val, isOK := storing.Get(msg.Key)
			if !isOK {
				return
			}

			response := transport.CreateGetResponseMessage(
				msg.RequestID,
				true,
				val,
			)

			peerToReply, ok := me.Cluster.GetPeer(
				msg.NodeID,
			)
			if !ok {
				return
			}

			err := transport.SendMessage(
				peerToReply,
				response,
			)
			if err != nil {
				return
			}
		} else {
			reject(msg.RequestID, me)
		}
	} else {
		// Here a new request needs to be made
		reject(msg.RequestID, me)
	}
}

func handleGetResponse(
	msg transport.ParsedMessage,
	me *Node,
) {
	req, ok := me.GetPendingRequest(
		msg.RequestID,
	)

	if !ok {
		return
	}

	req.ResultChan <- msg

	me.RemovePendingRequest(
		msg.RequestID,
	)
}

func handleGetREJ(
	msg transport.ParsedMessage,
	me *Node,
) {
	req, ok := me.GetPendingRequest(
		msg.RequestID,
	)

	if !ok {
		return
	}

	var finalResponse transport.ParsedMessage

	if msg.MembershipVersion > me.MembershipVersion {
		go SyncMembership(
			msg.NodeID,
			me,
			20*time.Second,
		)
	}

	if msg.SuggestedOwner != "" {

		newPeerToTry := cluster.NewPeer(
			msg.SuggestedOwner,
			msg.SuggestedAddr,
		)

		if !(me.Cluster.HasPeer(
			newPeerToTry.NodeID,
		)) {
			me.Cluster.AddPeer(
				newPeerToTry,
			)
		}

		thisRequest, ok := me.GetPendingRequest(
			msg.RequestID,
		)

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
			req.ResultChan <- finalResponse
			return
		}
		getResponseFromNewPeer, err := me.SendRequestAndWait(
			newPeerToTry,
			newMessage,
			30*time.Second,
		)
		if err != nil {

			finalResponse = msg
			req.ResultChan <- finalResponse
			return
		}

		finalResponse = getResponseFromNewPeer

	} else {
		finalResponse = msg
	}

	req.ResultChan <- finalResponse
}
