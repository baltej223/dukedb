package node

import (
	"fmt"
	"time"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/transport"
	"github.com/baltej223/dukedb/scripts"
)

func handleMembership(msg transport.ParsedMessage, me *Node) {
	currentPeers := make(map[string]string)
	for _, p := range me.Cluster.GetPeers() {
		currentPeers[p.NodeID] = p.Addr
	}

	for _, p := range msg.Peers {
		if _, exists := currentPeers[p.NodeID]; !exists {
			if p.NodeID != me.ID {
				me.Cluster.AddPeer(cluster.NewPeer(p.NodeID, p.Addr))
				me.MembershipVersion++
			}
		}
	}
}

func (me *Node) StartGossipLoop(printit bool) error {
	for {
		currentNeighbours := me.Cluster.GetPeers()
		if len(currentNeighbours) == 0 {
			time.Sleep(me.GossipLoopTime)
			continue
		}

		// Number of peers to send message to.
		numberOfPeers, err := scripts.RandomNumber(len(currentNeighbours))
		if err != nil {
			return err
		}

		randomPeers, err := scripts.ChooseRandomElements[cluster.Peer](
			currentNeighbours,
			numberOfPeers,
		)
		if err != nil {
			return err
		}

		gossipMessage, err := transport.CreateMembershipMessage(
			currentNeighbours,
			me.MembershipVersion,
		)
		if err != nil {
			return err
		}

		for _, target := range randomPeers {
			if target.NodeID == me.ID {
				continue
			}
			err := transport.SendMessage(target, gossipMessage)
			// log.Printf(
			// 	"[node=%s] gossiped membership (%d peers) to %s",
			// 	me.ID,
			// 	len(randomPeers),
			// 	target.NodeID,
			// )
			if err != nil {
				return err
			}
		}

		if printit {
			fmt.Println(me.Cluster.Dump())
		}

		time.Sleep(me.GossipLoopTime)
	}
}

func handleSYNCMembership(msg transport.ParsedMessage, me *Node) {
	gosspipMessage, err := transport.
		CreateSYNCMebershipResponseMessage(
			me.Cluster.GetPeers(),
			me.MembershipVersion,
			msg.RequestID)
	if err != nil {
		return
	}
	nodeToRespond, ok := me.Cluster.GetPeer(msg.NodeID)
	if !ok {
		return
	}
	err = transport.SendMessage(nodeToRespond, gosspipMessage)
	if err != nil {
		return
	}
}

func handleSYNCMembershipResponse(msg transport.ParsedMessage, me *Node) {
	PendingRequest, ok := me.GetPendingRequest(msg.RequestID)
	if !ok {
		return
	}
	handleMembership(msg, me)
	PendingRequest.ResultChan <- msg
}

func SyncMembership(nodeID string, me *Node, timeout time.Duration) {
	syncRequest, err := transport.CreateSYNCMembershipMessage(me.ID)
	if err != nil {
		return
	}
	peerToSendMessage, ok := me.Cluster.GetPeer(nodeID)
	if !ok {
		return
	}

	me.SendRequestAndWait(peerToSendMessage, syncRequest, timeout)
}
