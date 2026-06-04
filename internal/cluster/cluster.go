// Package cluster takes care of the cluster itself, knows about cluster, gossip lives here.
package cluster

import (
	"fmt"
	"strings"
	"sync"
)

type ClusterManager struct {
	mu sync.RWMutex

	Neighbours map[string]Peer
}

func NewClusterManager(
	peers []Peer,
) *ClusterManager {
	neighbours := make(
		map[string]Peer,
	)

	for _, peer := range peers {
		neighbours[peer.NodeID] = peer
	}

	return &ClusterManager{
		Neighbours: neighbours,
	}
}

func (c *ClusterManager) AddPeer(
	peer Peer,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Neighbours[peer.NodeID] = peer
}

func (c *ClusterManager) RemovePeer(
	nodeID string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.Neighbours, nodeID)
}

func (c *ClusterManager) GetPeer(
	nodeID string,
) (Peer, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	peer, ok := c.Neighbours[nodeID]
	return peer, ok
}

func (c *ClusterManager) HasPeer(
	nodeID string,
) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.Neighbours[nodeID]
	return ok
}

func (c *ClusterManager) GetPeers() []Peer {
	c.mu.RLock()
	defer c.mu.RUnlock()

	peers := make(
		[]Peer,
		0,
		len(c.Neighbours),
	)

	for _, peer := range c.Neighbours {
		peers = append(
			peers,
			peer,
		)
	}

	return peers
}

func (c *ClusterManager) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.Neighbours)
}

func (c *ClusterManager) Dump() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var b strings.Builder

	b.WriteString("Cluster Membership:\n")

	for _, peer := range c.Neighbours {
		b.WriteString(
			fmt.Sprintf(
				"  %s -> %s\n",
				peer.NodeID,
				peer.Addr,
			),
		)
	}

	return b.String()
}
