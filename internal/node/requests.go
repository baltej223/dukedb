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
	err := transport.SendMessage(peer, msg)
	if err != nil {
		return transport.ParsedMessage{}, err
	}

	response, err := n.WaitForPendingRequest(msg.RequestID, timeout)
	if err != nil {
		if errors.Is(err, ErrRequestTimedOut) {
			n.RemovePendingRequest(msg.RequestID)
			n.AddSuspectedDeadPeer(peer)
		}

		return transport.ParsedMessage{}, err
	}

	n.RemovePendingRequest(msg.RequestID)
	return response, nil
}
