// Package node owns the node
package node

type Node struct {
	ID       string
	hostname string

	// Cluster     *cluster.Manager
	// Transport 	 *transport.Transport
	// Storage     *storage.Engine
	// Router      *routing.Router
	// Replication *replication.Manager
}

func Initialise(ID string, hostname string) *Node {
	return &Node{
		ID,
		hostname,
	}
}

//
// func (n *Node) GetPort() string {
// 	return n.port
// }
//
// func (n *Node) GetFullHostname() string {
// 	return n.uri
// }
