package node

import (
	"log"

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

	log.Printf(
		"[node=%s] PUT calling FindOwner(%s)",
		me.ID,
		key,
	)

	keyOwner := routing.FindOwner(key, me.AllNodesSort())

	log.Printf(
		"[node=%s] PUT owner=%s self=%s",
		me.ID,
		keyOwner.NodeID,
		me.ID,
	)

	if keyOwner.NodeID == me.ID {
		log.Printf(
			"[node=%s] PUT key=%s sender=%s owner=%s",
			me.ID,
			key,
			msg.NodeID,
			keyOwner.NodeID,
		)

		storing.Put(key, msg.Value)

		log.Printf(
			"[node=%s] PUT stored key=%s",
			me.ID,
			key,
		)

		response := transport.CreatePutACKMessage(msg.RequestID)

		log.Printf(
			"[node=%s] PUT_ACK created request_id=%s",
			me.ID,
			msg.RequestID,
		)

		peerToReply, ok := me.Cluster.GetPeer(msg.NodeID)
		if !ok {
			log.Printf(
				"[node=%s] PUT_ACK failed: sender %q not found in cluster",
				me.ID,
				msg.NodeID,
			)
			return
		}

		log.Printf(
			"[node=%s] PUT_ACK sending to node=%s addr=%s",
			me.ID,
			peerToReply.NodeID,
			peerToReply.Addr,
		)

		err := transport.SendMessage(peerToReply, response)
		if err != nil {
			log.Printf(
				"[node=%s] PUT_ACK send failed: %v",
				me.ID,
				err,
			)
			return
		}

		log.Printf(
			"[node=%s] PUT_ACK sent successfully",
			me.ID,
		)

	} else {
		log.Printf(
			"[node=%s] PUT rejected owner=%s self=%s",
			me.ID,
			keyOwner.NodeID,
			me.ID,
		)

		response := transport.CreatePutREJMessage(msg.RequestID)

		peerToReply, ok := me.Cluster.GetPeer(msg.NodeID)
		if !ok {
			log.Printf(
				"[node=%s] PUT_REJ failed: sender %q not found in cluster",
				me.ID,
				msg.NodeID,
			)
			return
		}

		err := transport.SendMessage(peerToReply, response)
		if err != nil {
			log.Printf(
				"[node=%s] PUT_REJ send failed: %v",
				me.ID,
				err,
			)
			return
		}

		log.Printf(
			"[node=%s] PUT_REJ sent successfully",
			me.ID,
		)
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
	req.ResultChan <- msg
	me.RemovePendingRequest(msg.RequestID)
}
