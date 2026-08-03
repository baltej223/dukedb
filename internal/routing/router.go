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
	replicationFactor int,
) []cluster.Peer {
	if len(sortedNodes) <= 1 || replicationFactor <= 1 {
		return nil
	}

	h := fnv.New32a()
	h.Write([]byte(key))

	ownerIdx := int(h.Sum32()) % len(sortedNodes)

	// Total replicas needed = total copies - owner
	numReplicas := replicationFactor - 1

	// Can't have more replicas than other live nodes
	if numReplicas > len(sortedNodes)-1 {
		numReplicas = len(sortedNodes) - 1
	}

	replicas := make([]cluster.Peer, 0, numReplicas)

	for i := 1; i <= numReplicas; i++ {
		idx := (ownerIdx + i) % len(sortedNodes)
		replicas = append(replicas, sortedNodes[idx])
	}

	return replicas
}
