package node

import (
	"errors"
	"time"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/transport"
)

var ErrRequestTimedOut = errors.New(
	"request timed out",
)

func (n *Node) SendRequestAndWait(
	peer cluster.Peer,
	msg transport.Message,
	timeout time.Duration,
) (transport.ParsedMessage, error) {
	pendingRequest := PendingRequest{
		time.Now(),
		msg.Type,
		make(chan transport.ParsedMessage),
	}

	n.AddPendingRequest(msg.RequestID, &pendingRequest)
	transport.SendMessage(peer, msg)

	response, err := n.WaitForPendingRequest(msg.RequestID, timeout)
	if err != nil {
		if err == ErrRequestTimedOut {
			n.RemovePendingRequest(msg.RequestID)
			// Peer suspected to be unhealthy
		}
	}
	n.RemovePendingRequest(msg.RequestID)
	return response, nil
}
