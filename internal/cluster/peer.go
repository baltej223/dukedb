package cluster

type Peer struct {
	NodeID string
	Addr   string
}

func NewPeer(nodeid, addr string) Peer {
	return Peer{nodeid, addr}
}
