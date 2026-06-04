// Package node owns the node
package node

import (
	"cmp"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/transport"
)

type Node struct {
	ID       string
	Hostname string

	PendingRequests map[string]*PendingRequest
	PendingMu       sync.RWMutex

	SuspectedDeadPeers   SuspectedDeadPeers
	SuspectedDeadPeersMu sync.RWMutex
	Cluster              *cluster.ClusterManager

	GossipLoopTime time.Duration
	// Transport 	 *transport.Transport
	// Storage     *storage.Engine
	// Router      *routing.Router
	// Replication *replication.Manager
}

type PendingRequest struct {
	CreatedAt time.Time
	Type      transport.MessageType

	ResultChan chan transport.ParsedMessage
}

func Initialise(
	ID string,
	hostname string,
	peers []cluster.Peer,
	GossipLoopTime time.Duration,
) *Node {
	return &Node{
		ID:       ID,
		Hostname: hostname,

		PendingRequests: make(
			map[string]*PendingRequest,
		),

		SuspectedDeadPeers: make(
			SuspectedDeadPeers,
			0,
		),

		Cluster: cluster.NewClusterManager(
			peers,
		),

		GossipLoopTime: GossipLoopTime,
	}
}

func (n *Node) AddPendingRequest(
	requestID string,
	req *PendingRequest,
) {
	n.PendingMu.Lock()
	defer n.PendingMu.Unlock()

	n.PendingRequests[requestID] = req
}

func (n *Node) RemovePendingRequest(requestID string) {
	n.PendingMu.Lock()
	defer n.PendingMu.Unlock()

	delete(n.PendingRequests, requestID)
}

func (n *Node) GetPendingRequest(requestID string) (*PendingRequest, bool) {
	n.PendingMu.RLock()
	defer n.PendingMu.RUnlock()

	req, ok := n.PendingRequests[requestID]
	return req, ok
}

func (n *Node) WaitForPendingRequest(
	requestID string,
	timeout time.Duration,
) (transport.ParsedMessage, error) {
	req, ok := n.GetPendingRequest(requestID)
	if !ok {
		return transport.ParsedMessage{},
			fmt.Errorf(
				"pending request %s not found",
				requestID,
			)
	}

	select {

	case response := <-req.ResultChan:
		return response, nil

	case <-time.After(timeout):
		return transport.ParsedMessage{},
			ErrRequestTimedOut
	}
}

func (me *Node) GetAllNodes() []cluster.Peer {
	neighbours := me.Cluster.GetPeers()
	selfAsPeer := cluster.NewPeer(me.ID, me.Hostname)
	return append(neighbours, selfAsPeer)
}

func (me *Node) AllNodesSort() []cluster.Peer {
	currentNodes := me.GetAllNodes()
	slices.SortFunc(currentNodes, func(a, b cluster.Peer) int {
		return cmp.Compare(a.NodeID, b.NodeID)
	})
	return currentNodes
}
