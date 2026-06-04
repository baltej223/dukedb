package routing

import (
	"hash/fnv"

	"github.com/baltej223/dukedb/internal/cluster"
)

type Router struct{}

func FindOwner(
	key string,
	sortedNodes []cluster.Peer, //
) cluster.Peer {
	h := fnv.New32a()

	h.Write([]byte(key))

	value := h.Sum32()

	idx := int(value) % len(sortedNodes)
	return sortedNodes[idx]
}
