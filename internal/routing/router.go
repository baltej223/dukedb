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

func FindReplicas(
	key string,
	sortedNodes []cluster.Peer,
	ReplicationFactor int,
) []cluster.Peer {
	replicas := []cluster.Peer{}

	h := fnv.New32a()
	h.Write([]byte(key))

	value := h.Sum32()
	for i := 1; i <= ReplicationFactor; i++ {
		idx := (int(value) + i) % len(sortedNodes)
		replicas = append(replicas, sortedNodes[idx])
	}
	return replicas
}
