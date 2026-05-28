package node

type Node struct {
	ID   string
	host string
	port string
	uri  string
	// Cluster     *cluster.Manager
	// Transport 	 *transport.Transport
	// Storage     *storage.Engine
	// Router      *routing.Router
	// Replication *replication.Manager
}

func Initialise(ID string, host string, port string) *Node {
	return &Node{
		ID,
		host,
		port,
		host + ":" + port,
	}
}

func (n *Node) GetPort() string {
	return n.port
}

func (n *Node) GetFullHostname() string {
	return n.uri
}
