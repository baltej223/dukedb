package node

import "github.com/baltej223/dukedb/internal/cluster"

type SuspectedDeadPeers = []cluster.Peer

// func (n *Node) SetSuspectedDeadPeers(
// 	peers SuspectedDeadPeers,
// ) {
// 	n.SuspectedDeadPeersMu.Lock()
// 	defer n.SuspectedDeadPeersMu.Unlock()
//
// 	n.SuspectedDeadPeers = peers
// }
//
// func (n *Node) GetSuspectedDeadPeers() SuspectedDeadPeers {
// 	n.SuspectedDeadPeersMu.RLock()
// 	defer n.SuspectedDeadPeersMu.RUnlock()
//
// 	return n.SuspectedDeadPeers
// }

func (n *Node) AddSuspectedDeadPeer(
	peer cluster.Peer,
) {
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
