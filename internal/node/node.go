// Package node owns the node
package node

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/routing"
	"github.com/baltej223/dukedb/internal/storing"
	"github.com/baltej223/dukedb/internal/transport"
	dukelog "github.com/baltej223/dukedb/log"
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

	MembershipVersion   int
	MembershipVersionMu sync.RWMutex

	ReplicationFactor int

	TransportServer *transport.Server

	// Transport 	 *transport.Transport
	// Storage     *storage.Engine
	// Router      *routing.Router
	// Replication *replication.Manager
}

type PendingRequest struct {
	CreatedAt time.Time

	Message    transport.Message
	RetryCount int

	ResultChan chan transport.ParsedMessage
}

func Initialise(
	ID string,
	hostname string,
	peers []cluster.Peer,
	GossipLoopTime time.Duration,
	ReplicationFactor int,
	TransportServer *transport.Server,
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

		GossipLoopTime:    GossipLoopTime,
		MembershipVersion: 0,

		ReplicationFactor: ReplicationFactor,

		TransportServer: TransportServer,
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

func (n *Node) RemovePendingRequest(
	requestID string,
) {
	n.PendingMu.Lock()
	defer n.PendingMu.Unlock()

	delete(
		n.PendingRequests,
		requestID,
	)
}

func (n *Node) GetPendingRequest(
	requestID string,
) (*PendingRequest, bool) {
	n.PendingMu.RLock()
	defer n.PendingMu.RUnlock()

	req, ok := n.PendingRequests[requestID]

	if !ok {
		dukelog.Printf(
			"[PENDING_MISS] request_id=%s",
			requestID,
		)
	}

	return req, ok
}

func (n *Node) WaitForPendingRequest(
	requestID string,
	timeout time.Duration,
) (transport.ParsedMessage, error) {
	req, ok := n.GetPendingRequest(
		requestID,
	)
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

type ErrorCode error

var (
	KeyNotExists ErrorCode = errors.New("KEY NOT EXISTS")
	ErrorOccured ErrorCode = errors.New("ERROR OCCURED")
	PutFailed    ErrorCode = errors.New("PUT FAILED")
)

func GET(key string, me *Node) (string, error) {
	owner := routing.FindOwner(key, me.AllNodesSort())
	dukelog.Printf("CLIENT GET REQ RECEIVED, owner node?: %s", owner.NodeID == me.ID)
	if owner.NodeID == me.ID {
		dukelog.Print("Own storing is used.")
		if storing.Exists(key) {
			value, ok := storing.Get(key)
			if !ok {
				return "", ErrorOccured
			}
			return string(value), nil
		} else {
			return "", KeyNotExists
		}
	} else {

		request, err := transport.CreateGetMessage(
			key,
			me.ID,
			me.MembershipVersion,
		)
		if err != nil {
			return "", err
		}
		response, err := me.SendRequestAndWait(owner, request, 10*time.Second)

		if err == nil && response.Found {
			return string(response.Value), nil
		}

		// Try replicas
		replicas := routing.FindReplicas(
			key,
			me.AllNodesSort(),
			me.ReplicationFactor,
		)

		for _, replica := range replicas {
			if replica.NodeID == me.ID {
				continue
			}
			response, err = me.SendRequestAndWait(
				replica,
				request,
				10*time.Second,
			)
			if err != nil {
				continue
			}

			if response.Found {
				return string(response.Value), nil
			}
		}

		return "", KeyNotExists
	}
}

func PUT(key, value string, me *Node) error {
	owner := routing.FindOwner(key, me.AllNodesSort())
	if owner.NodeID == me.ID {
		storing.Put(key, []byte(value))
		return nil
	} else {
		request, err := transport.CreatePutMessage(
			key,
			[]byte(value),
			me.ID,
			me.MembershipVersion,
		)
		if err != nil {
			return err
		}
		response, err := me.SendRequestAndWait(owner, request, 20*time.Second)
		if err != nil {
			return err
		}
		if response.Success {
			return nil
		} else {
			return PutFailed
		}
	}
}

func GetMembershipVersion(me *Node) int {
	me.MembershipVersionMu.RLock()
	defer me.MembershipVersionMu.RUnlock()

	return me.MembershipVersion
}

func IncreaseMembershipVersion(me *Node) {
	me.MembershipVersionMu.Lock()
	defer me.MembershipVersionMu.Unlock()

	me.MembershipVersion = me.MembershipVersion + 1
}

func UpdateMembershipVersion(me *Node, newVersion int) {
	me.MembershipVersionMu.Lock()
	defer me.MembershipVersionMu.Unlock()

	me.MembershipVersion = newVersion
}
