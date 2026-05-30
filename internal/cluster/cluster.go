// Package cluster takes care of the cluster itself, knows about cluster, gossip lives here.
package cluster

type Peer struct {
	NodeID string
	Addr   string
}
