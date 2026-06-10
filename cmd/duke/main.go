package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/baltej223/dukedb/internal/api"
	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/node"
	"github.com/baltej223/dukedb/internal/storing"
	"github.com/baltej223/dukedb/internal/transport"
)

func main() {
	// Flags handling
	selfAddress := flag.String("self-addr", "localhost:8000", "Address of the current node, Example localhost:8000")
	selfNodeID := flag.String("self-node-id", "-", "Self node ID, Example: b")

	isSeedNode := flag.Bool("seed-node", false, "Is this a seed node.")

	peerAddress := flag.String("peer-addr", "", "Address of peer node, Example: localhost:8001")
	peerNodeID := flag.String("peer-node-id", "", "Peer node ID, Example: b")
	delay := flag.Int("delay", 5, "[Debug]: Initial Delay Before sending first request")
	apiAt := flag.String("api-at", ":9000", "Where to run API server at?")
	// FLAGS END
	flag.Parse()

	// Flags check
	if *isSeedNode {
		if *peerAddress != "" || *peerNodeID != "" {
			panic("Peers can't be defined for seed node.")
		}
	} else {
		if *peerAddress == "" || *peerNodeID == "" {
			panic("One Peer should be defined for a non seed node.")
		}
	}
	// Flags Check END

	hostname := *selfAddress
	var me *node.Node
	var neighbours []cluster.Peer
	if *isSeedNode {
		neighbours = []cluster.Peer{}
	} else {
		firstPeer := cluster.NewPeer(*peerNodeID, *peerAddress)
		neighbours = []cluster.Peer{firstPeer}
	}

	me = node.Initialise(
		*selfNodeID,
		*selfAddress,
		neighbours,
		10*time.Second,
	)
	storing.InitialiseKV()

	// Init tranport server
	server := transport.NewServer(hostname)
	log.Println("Starting duke node on " + me.Hostname)

	go func() {
		err := server.Start(func(conn net.Conn) {
			transport.HandleConnection(
				conn,
				func(msg transport.ParsedMessage) {
					node.Dispatch(msg, me)
				},
			)
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done")
	}()
	// END

	// Needed
	if *peerAddress != "" && !*isSeedNode {
		time.Sleep(time.Duration(*delay) * time.Second)
		joingRequest,
			err := transport.CreateJoinMessage(
			*selfNodeID,
			*selfAddress,
		)
		if err != nil {
			panic(err)
		}

		_, err = me.SendRequestAndWait(
			cluster.NewPeer(*peerNodeID, *peerAddress),
			joingRequest,
			100*time.Second,
		)
		if err != nil {
			panic(err)
		}
	}

	// Gossip loop
	go func() {
		err := me.StartGossipLoop(false)
		if err != nil {
			log.Printf("gossip failed: %v", err)
		}
	}()

	go func() {
		apiServer := api.NewServer(*apiAt, me)
		apiServer.Start()
	}()

	select {}
}
