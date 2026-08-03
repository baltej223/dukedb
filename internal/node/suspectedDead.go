package node

import (
	"errors"
	"time"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/transport"
	dukelog "github.com/baltej223/dukedb/log"
)

type SuspectedDeadPeers = []cluster.Peer

func (n *Node) AddSuspectedDeadPeer(
	peer cluster.Peer,
) {
	if peer.NodeID == n.ID {
		panic("trying to suspect myself")
	}
	n.SuspectedDeadPeersMu.Lock()
	defer n.SuspectedDeadPeersMu.Unlock()

	for _, p := range n.SuspectedDeadPeers {
		if p.NodeID == peer.NodeID {
			return
		}
	}

	n.SuspectedDeadPeers = append(
		n.SuspectedDeadPeers,
		peer,
	)
}

func (n *Node) RemoveSuspectedDeadPeer(
	nodeID string,
) {
	n.SuspectedDeadPeersMu.Lock()
	defer n.SuspectedDeadPeersMu.Unlock()

	filtered := make(
		SuspectedDeadPeers,
		0,
		len(n.SuspectedDeadPeers),
	)

	for _, peer := range n.SuspectedDeadPeers {
		if peer.NodeID != nodeID {
			filtered = append(
				filtered,
				peer,
			)
		}
	}

	n.SuspectedDeadPeers = filtered
}

func (n *Node) IsSuspectedDead(
	nodeID string,
) bool {
	n.SuspectedDeadPeersMu.RLock()
	defer n.SuspectedDeadPeersMu.RUnlock()

	for _, peer := range n.SuspectedDeadPeers {
		if peer.NodeID == nodeID {
			return true
		}
	}

	return false
}

func (me *Node) GetSuspectedPeers() SuspectedDeadPeers {
	me.SuspectedDeadPeersMu.RLock()
	defer me.SuspectedDeadPeersMu.RUnlock()

	return me.SuspectedDeadPeers
}

func (me *Node) StartSuspectedPeerChecker(tickerTimer time.Duration) {
	ticker := time.NewTicker(tickerTimer)
	defer ticker.Stop()

	for range ticker.C {
		me.CheckSuspectedPeers()
	}
}

func (me *Node) CheckSuspectedPeers() {
	peers := me.GetSuspectedPeers()
	dukelog.Printf("Suspecting peers, Number of deadpeers: %d", len(peers))
	for _, peer := range peers {
		pingMessage, err := transport.CreatePingMessage(peer.NodeID)
		if err != nil {
			dukelog.Error(err)
		}
		_, err = me.SendRequestAndWaitWithoutDeadCheckPeer(peer, pingMessage, 30*time.Second)
		if err != nil {
			if errors.Is(
				err,
				ErrRequestTimedOut,
			) {
				// Is SuspectedDeadPeer
				// remove from cluster
				// membership++
				// gossip
				dukelog.Printf("A suspecteddeadpeer found: %s", peer.NodeID)
				me.Cluster.RemovePeer(peer.NodeID)
				me.MembershipVersionMu.Lock()
				defer me.MembershipVersionMu.Unlock()
				me.MembershipVersion++
			} else {
				me.RemoveSuspectedDeadPeer(peer.NodeID)
			}
		} else {
			dukelog.Errorf("Some error occured: %s", err)
		}
	}
}
