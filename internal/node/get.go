package node

import (
	"log"
	"time"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/routing"
	"github.com/baltej223/dukedb/internal/storing"
	"github.com/baltej223/dukedb/internal/transport"
)

// func handleGet(msg transport.ParsedMessage, me *Node) {
// 	sortedNodes := me.AllNodesSort()
// 	ownerNode := routing.FindOwner(msg.Key, sortedNodes)
//
// 	if ownerNode.NodeID == me.ID {
// 		if storing.Exists(msg.Key) {
// 			val, isOK := storing.Get(msg.Key)
// 			if !isOK {
// 				return
// 			}
// 			response := transport.CreateGetResponseMessage(msg.RequestID, true, val)
//
// 			peerToReply, ok := me.Cluster.GetPeer(msg.NodeID)
// 			if !ok {
// 				return
// 			}
// 			err := transport.SendMessage(peerToReply, response)
// 			if err != nil {
// 				return
// 			}
// 		}
// 	} else {
// 		log.Printf(
// 			"KEY MISSING key=%s",
// 			msg.Key,
// 		)
// 		// Here a new request needs to be made
// 	}
// }
//
// func handleGetResponse(msg transport.ParsedMessage, me *Node) {
// 	req, ok := me.GetPendingRequest(msg.RequestID)
// 	if !ok {
// 		log.Printf(
// 			"[node=%s] no pending request found for request_id=%s",
// 			me.ID,
// 			msg.RequestID,
// 		)
// 		return
// 	}
// 	req.ResultChan <- msg
//
// 	log.Printf(
// 		"[node=%s] pending request fulfilled request_id=%s",
// 		me.ID,
// 		msg.RequestID,
// 	)
// }
//
// func handleGetREJ(msg transport.ParsedMessage, me *Node) {
// 	req, ok := me.GetPendingRequest(msg.RequestID)
// 	if !ok {
// 		return
// 	}
// 	var finalResponse transport.ParsedMessage
//
// 	if msg.MembershipVersion > me.MembershipVersion {
// 		// schedule membership sync
// 		go SyncMembership(msg.NodeID, me, 20*time.Second)
// 	}
// 	if msg.SuggestedOwner != "" {
// 		newPeerToTry := cluster.NewPeer(msg.SuggestedOwner, msg.SuggestedAddr)
// 		if !(me.Cluster.HasPeer(newPeerToTry.NodeID)) {
// 			me.Cluster.AddPeer(newPeerToTry)
// 		}
//
// 		thisRequest, ok := me.GetPendingRequest(msg.RequestID)
// 		if !ok {
// 			return
// 		}
// 		newMessage, err := transport.CreateGetMessage(
// 			thisRequest.Message.Headers["KEY"],
// 			me.ID,
// 			me.MembershipVersion,
// 		)
// 		if err != nil {
// 			finalResponse = msg
// 			req.ResultChan <- finalResponse
// 			return
// 		}
//
// 		getResponseFromNewPeer, err := me.SendRequestAndWait(newPeerToTry, newMessage, 30*time.Second)
// 		if err != nil {
// 			finalResponse = msg
// 			req.ResultChan <- finalResponse
// 			return
// 		}
// 		finalResponse = getResponseFromNewPeer
// 	} else {
// 		finalResponse = msg
// 	}
//
// 	req.ResultChan <- finalResponse
// }

func handleGet(
	msg transport.ParsedMessage,
	me *Node,
) {
	log.Printf(
		"[GET_ENTER] request_id=%s key=%s node=%s sender=%s",
		msg.RequestID,
		msg.Key,
		me.ID,
		msg.NodeID,
	)

	sortedNodes := me.AllNodesSort()

	ownerNode := routing.FindOwner(
		msg.Key,
		sortedNodes,
	)

	log.Printf(
		"[ROUTE] op=GET key=%s owner=%s ring=%v",
		msg.Key,
		ownerNode.NodeID,
		me.AllNodesSort(),
	)

	log.Printf(
		"[GET_OWNER] request_id=%s key=%s owner=%s me=%s",
		msg.RequestID,
		msg.Key,
		ownerNode.NodeID,
		me.ID,
	)

	if ownerNode.NodeID == me.ID {

		log.Printf(
			"[GET_LOCAL_OWNER] request_id=%s key=%s",
			msg.RequestID,
			msg.Key,
		)
		exists := storing.Exists(msg.Key)

		log.Printf(
			"[GET_EXISTS_CHECK] request_id=%s key=%s exists=%v",
			msg.RequestID,
			msg.Key,
			exists,
		)
		if storing.Exists(msg.Key) {

			log.Printf(
				"[GET_EXISTS] request_id=%s key=%s",
				msg.RequestID,
				msg.Key,
			)

			val, isOK := storing.Get(msg.Key)
			if !isOK {

				log.Printf(
					"[GET_GET_FAILED] request_id=%s key=%s",
					msg.RequestID,
					msg.Key,
				)

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

				log.Printf(
					"[GET_REPLY_PEER_NOT_FOUND] request_id=%s sender=%s",
					msg.RequestID,
					msg.NodeID,
				)

				return
			}

			log.Printf(
				"[GET_SEND_RESPONSE] request_id=%s target=%s",
				msg.RequestID,
				peerToReply.NodeID,
			)

			err := transport.SendMessage(
				peerToReply,
				response,
			)
			if err != nil {

				log.Printf(
					"[GET_SEND_RESPONSE_ERR] request_id=%s err=%v",
					msg.RequestID,
					err,
				)

				return
			}
		}
	} else {
		log.Printf(
			"[GET_NOT_OWNER] request_id=%s key=%s owner=%s me=%s",
			msg.RequestID,
			msg.Key,
			ownerNode.NodeID,
			me.ID,
		)

		// Here a new request needs to be made
	}
}

func handleGetResponse(
	msg transport.ParsedMessage,
	me *Node,
) {
	log.Printf(
		"[GET_RESPONSE_ENTER] request_id=%s found=%v",
		msg.RequestID,
		msg.Found,
	)

	req, ok := me.GetPendingRequest(
		msg.RequestID,
	)

	if !ok {

		log.Printf(
			"[GET_RESPONSE_ORPHAN] request_id=%s",
			msg.RequestID,
		)

		return
	}

	log.Printf(
		"[GET_RESPONSE_FULFILL] request_id=%s",
		msg.RequestID,
	)

	req.ResultChan <- msg

	log.Printf(
		"[GET_RESPONSE_FULFILLED] request_id=%s",
		msg.RequestID,
	)

	me.RemovePendingRequest(
		msg.RequestID,
	)

	log.Printf(
		"[GET_RESPONSE_REMOVED_PENDING] request_id=%s",
		msg.RequestID,
	)
}

func handleGetREJ(
	msg transport.ParsedMessage,
	me *Node,
) {
	log.Printf(
		"[GET_REJ_ENTER] request_id=%s suggested_owner=%s suggested_addr=%s",
		msg.RequestID,
		msg.SuggestedOwner,
		msg.SuggestedAddr,
	)

	req, ok := me.GetPendingRequest(
		msg.RequestID,
	)

	if !ok {

		log.Printf(
			"[GET_REJ_PENDING_NOT_FOUND] request_id=%s",
			msg.RequestID,
		)

		return
	}

	var finalResponse transport.ParsedMessage

	if msg.MembershipVersion > me.MembershipVersion {

		log.Printf(
			"[GET_REJ_SYNC_MEMBERSHIP] request_id=%s local=%d remote=%d",
			msg.RequestID,
			me.MembershipVersion,
			msg.MembershipVersion,
		)

		go SyncMembership(
			msg.NodeID,
			me,
			20*time.Second,
		)
	}

	if msg.SuggestedOwner != "" {

		log.Printf(
			"[GET_REJ_RETRY_START] request_id=%s owner=%s",
			msg.RequestID,
			msg.SuggestedOwner,
		)

		newPeerToTry := cluster.NewPeer(
			msg.SuggestedOwner,
			msg.SuggestedAddr,
		)

		if !(me.Cluster.HasPeer(
			newPeerToTry.NodeID,
		)) {

			log.Printf(
				"[GET_REJ_ADD_PEER] request_id=%s peer=%s",
				msg.RequestID,
				newPeerToTry.NodeID,
			)

			me.Cluster.AddPeer(
				newPeerToTry,
			)
		}

		thisRequest, ok := me.GetPendingRequest(
			msg.RequestID,
		)

		if !ok {

			log.Printf(
				"[GET_REJ_REQUEST_GONE] request_id=%s",
				msg.RequestID,
			)

			return
		}

		newMessage, err := transport.CreateGetMessage(
			thisRequest.Message.Headers["KEY"],
			me.ID,
			me.MembershipVersion,
		)
		if err != nil {

			log.Printf(
				"[GET_REJ_BUILD_RETRY_FAILED] request_id=%s err=%v",
				msg.RequestID,
				err,
			)

			finalResponse = msg
			req.ResultChan <- finalResponse
			return
		}

		log.Printf(
			"[GET_REJ_RETRY_SEND] request_id=%s target=%s",
			msg.RequestID,
			newPeerToTry.NodeID,
		)

		getResponseFromNewPeer, err := me.SendRequestAndWait(
			newPeerToTry,
			newMessage,
			30*time.Second,
		)
		if err != nil {

			log.Printf(
				"[GET_REJ_RETRY_FAILED] request_id=%s err=%v",
				msg.RequestID,
				err,
			)

			finalResponse = msg
			req.ResultChan <- finalResponse
			return
		}

		log.Printf(
			"[GET_REJ_RETRY_SUCCESS] request_id=%s",
			msg.RequestID,
		)

		finalResponse = getResponseFromNewPeer

	} else {

		log.Printf(
			"[GET_REJ_NO_SUGGESTED_OWNER] request_id=%s",
			msg.RequestID,
		)

		finalResponse = msg
	}

	log.Printf(
		"[GET_REJ_FINAL_RESPONSE] request_id=%s",
		msg.RequestID,
	)

	req.ResultChan <- finalResponse
}
