package node

import (
	"errors"
	"time"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/transport"
	dukelog "github.com/baltej223/dukedb/log"
)

var ErrRequestTimedOut = errors.New(
	"request timed out",
)

func (me *Node) SendRequestAndWait(
	peer cluster.Peer,
	msg transport.Message,
	timeout time.Duration,
) (transport.ParsedMessage, error) {
	if peer.NodeID == me.ID {
		panic("SendRequestAndWait called on myself")
	}

	defer me.RemovePendingRequest(
		msg.RequestID,
	)

	pendingRequest := PendingRequest{
		CreatedAt: time.Now(),
		Message:   msg,
		ResultChan: make(
			chan transport.ParsedMessage,
		),
	}

	me.AddPendingRequest(
		msg.RequestID,
		&pendingRequest,
	)

	err := transport.SendMessage(
		peer,
		msg,
	)
	if err != nil {
		return transport.ParsedMessage{}, err
	}

	response, err := me.WaitForPendingRequest(
		msg.RequestID,
		timeout,
	)
	if err != nil {

		if errors.Is(
			err,
			ErrRequestTimedOut,
		) && !me.IsSuspectedDead(peer.NodeID) {

			pingMessage, err := transport.CreatePingMessage(me.ID)
			if err != nil {
				dukelog.Error(err)
			}
			me.AddSuspectedDeadPeer(
				peer,
			)
			err = transport.SendMessage(peer, pingMessage)
			if err != nil {
				dukelog.Error("Error in sending ping for node testing.")
			}
		}

		// Else it can also mean a network partition.
		// else: will be hadled by the deadcheck loop
		// Branches can be set here to determine what the issue is, and handlers can be wrritten for them accordingly

		return transport.ParsedMessage{},
			err
	}
	return response, nil
}

func (me *Node) SendRequestAndWaitWithoutDeadCheckPeer(
	peer cluster.Peer,
	msg transport.Message,
	timeout time.Duration,
) (transport.ParsedMessage, error) {
	if peer.NodeID == me.ID {
		panic("SendRequestAndWait called on myself")
	}

	defer me.RemovePendingRequest(
		msg.RequestID,
	)

	pendingRequest := PendingRequest{
		CreatedAt: time.Now(),
		Message:   msg,
		ResultChan: make(
			chan transport.ParsedMessage,
		),
	}

	me.AddPendingRequest(
		msg.RequestID,
		&pendingRequest,
	)

	err := transport.SendMessage(
		peer,
		msg,
	)
	if err != nil {
		return transport.ParsedMessage{}, err
	}

	response, _ := me.WaitForPendingRequest(
		msg.RequestID,
		timeout,
	)
	return response, nil
}
