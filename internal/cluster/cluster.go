// Package cluster takes care of the cluster itself, knows about cluster, gossip lives here.
package cluster

import (
	"errors"

	"github.com/baltej223/dukedb/internal/storing"
)

type Peer struct {
	NodeID string
	Addr   string
}

// get peer from node id

func PeerFromNodeID(NodeIDOfPeer string) (Peer, error) {
	neighbours, _ := storing.GetIJSON[[]Peer]("neighbours")

	for i := 0; i < len(neighbours); i++ {
		if neighbours[i].NodeID == NodeIDOfPeer {
			return neighbours[i], nil
		}
	}

	return Peer{}, errors.New("Error")
}
