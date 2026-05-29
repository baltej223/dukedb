// Package cluster takes care of the cluster itself, knows about cluster, gossip lives here.
package cluster

type OtherNode struct {
	Name    string
	Address string // Should literally be the address where the call will be made
}
