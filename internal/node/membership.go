package node

import (
	"errors"
	"time"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/dukerror"
	"github.com/baltej223/dukedb/internal/transport"
	dukelog "github.com/baltej223/dukedb/log"
	"github.com/baltej223/dukedb/scripts"
)

func handleMembership(msg transport.ParsedMessage, me *Node) {
	currentVersion := GetMembershipVersion(me)

	if msg.MembershipVersion < currentVersion {
		return
	}

	changed := false

	// Add peers that we don't already have.
	for _, p := range msg.Peers {
		if p.NodeID == me.ID {
			continue
		}

		if _, exists := me.Cluster.GetPeer(p.NodeID); exists {
			continue
		}

		me.Cluster.AddPeer(p)
		changed = true
	}

	if msg.MembershipVersion > currentVersion {
		UpdateMembershipVersion(me, msg.MembershipVersion)
		changed = true
	}

	// Merge suspected-dead information.
	for _, p := range msg.SuspectedDeadPeers {
		if p.Peer.NodeID == me.ID {
			continue
		}

		if me.IsSuspectedDead(p.Peer.NodeID) {
			continue
		}

		me.AddSuspectedDeadPeer(p.Peer)
		changed = true
	}

	// Local membership changes need a new version.
	if changed && msg.MembershipVersion <= currentVersion {
		IncreaseMembershipVersion(me)
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
			me.GetSuspectedPeers(),
			GetMembershipVersion(me),
		)
		if err != nil {
			return err
		}

		for _, target := range randomPeers {
			if target.NodeID == me.ID {
				continue
			}
			err := transport.SendMessage(target, gossipMessage)
			if err != nil {
				if errors.Is(dukerror.Normalize(err), dukerror.ErrNetwork) {
					me.AddSuspectedDeadPeer(target)
				}
			}
		}

		if printit {
			dukelog.Print(me.Cluster.Dump())
		}

		time.Sleep(me.GossipLoopTime)
	}
}

func handleSYNCMembership(msg transport.ParsedMessage, me *Node) {
	gosspipMessage, err := transport.
		CreateSYNCMebershipResponseMessage(
			me.Cluster.GetPeers(),
			GetMembershipVersion(me),
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
