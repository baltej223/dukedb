// Package node owns the node
package node

import (
	"fmt"
	"sync"
	"time"

	"github.com/baltej223/dukedb/internal/transport"
)

type Node struct {
	ID       string
	Hostname string

	PendingRequests map[string]*PendingRequest
	PendingMu       sync.RWMutex
	// Cluster     *cluster.Manager
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

func Initialise(ID string, hostname string) *Node {
	return &Node{
		ID:       ID,
		Hostname: hostname,

		PendingRequests: make(
			map[string]*PendingRequest,
		),
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
