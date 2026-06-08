package node

import (
	"fmt"
	"log"
	"time"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/transport"
	"github.com/baltej223/dukedb/scripts"
)

// [UGLY]: This function is very ugly, I need to write better version of this function in future/
func handleMembership(msg transport.ParsedMessage, me *Node) {
	// here I will compare the infomation which has been sent by the other node, and the infomration which I have.

	newPeers := make(map[string]string)
	currentPeers := make(map[string]string)

	// First creating the map of received infomation.

	for _, val := range msg.Peers {
		newPeers[val.NodeID] = val.Addr
	}

	for _, val := range me.Cluster.GetPeers() {
		currentPeers[val.NodeID] = val.Addr
	}

	for _, peer := range msg.Peers {
		if _, ok := currentPeers[peer.NodeID]; !ok {
			// The sent info not here.
			newPeer := cluster.NewPeer(peer.NodeID, newPeers[peer.NodeID])
			me.Cluster.AddPeer(newPeer)
			me.MembershipVersion = me.MembershipVersion + 1
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
			log.Printf(
				"[node=%s] gossiped membership (%d peers) to %s",
				me.ID,
				len(randomPeers),
				target.NodeID,
			)
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
