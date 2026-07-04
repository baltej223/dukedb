package node

import (
	"time"

	"github.com/baltej223/dukedb/internal/routing"
	"github.com/baltej223/dukedb/internal/transport"
	dukelog "github.com/baltej223/dukedb/log"
)

func SendPutToReplicas(key string, value string, me *Node) {
	replicas := routing.FindReplicas(
		key,
		me.AllNodesSort(),
		me.ReplicationFactor,
	)

	for _, replica := range replicas {
		if me.IsSuspectedDead(replica.NodeID) {
			continue
		}
		put, err := transport.CreatePutMessage(
			key,
			[]byte(value),
			me.ID,
			me.MembershipVersion,
		)
		if err != nil {
			dukelog.Error("Some Error Occured while creating a Put message for sending to replicas")
		}

		response, err := me.SendRequestAndWait(
			replica,
			put,
			60*time.Second,
		)
		if err != nil {
			dukelog.Error("Error while streaming the keys to replicas: " + err.Error())
		}
		if response.Success {
			dukelog.Printf("Successfully streamed keys to replicas.")
		}
	}
}

func IsReplica(key string, me *Node) bool {
	keyReplicas := routing.FindReplicas(
		key,
		me.AllNodesSort(),
		me.ReplicationFactor,
	)

	for _, keyReplica := range keyReplicas {
		if keyReplica.NodeID == me.ID {
			return true
		}
	}
	return false
}
