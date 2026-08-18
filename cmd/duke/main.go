package main

import (
	"flag"
	"net"
	"time"

	"github.com/baltej223/dukedb/internal/api"
	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/dukerror"
	"github.com/baltej223/dukedb/internal/node"
	"github.com/baltej223/dukedb/internal/storing"
	"github.com/baltej223/dukedb/internal/transport"
	dukelog "github.com/baltej223/dukedb/log"
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
	rf := flag.Int("replication-factor", 1, "The shared replication factor of the cluster.")
	timerTicker := flag.Int("dead-check-ticker-timer", 10, "The time to wait before macking an another check for all the suspected dead peers.")
	printMembershipInfo := flag.Bool("print-membership-info", false, "<DEBUG> Should the node print the cluster membership info.")
	// FLAGS END
	flag.Parse()

	// Flags check
	if *isSeedNode {
		if *peerAddress != "" || *peerNodeID != "" {
			dukerror.Handle(dukerror.ErrPeersDefinedForSeedNode)
		}
	} else {
		if *peerAddress == "" || *peerNodeID == "" {
			dukerror.Handle(dukerror.ErrPeersDefinedForSeedNode)
		}
	}
	// Flags Check END

	// Start Logging here
	defer dukelog.Flush()
	// ------------------

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
		*rf,
		&transport.Server{},
	)
	storing.InitialiseKV()

	// Init tranport server
	server := transport.NewServer(hostname)
	dukelog.Printf("Starting duke node on %s", me.Hostname)

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
			dukerror.Handle(dukerror.Normalize(err))
		} else {
			me.TransportServer = server
		}
		dukelog.Println("Server: Done")
	}()
	// END

	// Initial Joining request
	if *peerAddress != "" && !*isSeedNode {
		time.Sleep(time.Duration(*delay) * time.Second)
		joingRequest,
			err := transport.CreateJoinMessage(
			*selfNodeID,
			*selfAddress,
		)
		if err != nil {
			dukerror.Handle(dukerror.Normalize(err))
		}

		_, err = me.SendRequestAndWait(
			cluster.NewPeer(*peerNodeID, *peerAddress),
			joingRequest,
			100*time.Second,
		)
		if err != nil {
			dukerror.Handle(dukerror.Normalize(dukerror.ErrClusterJoiningFailed))
		}
	}

	// Gossip loop
	go func() {
		err := me.StartGossipLoop(*printMembershipInfo)
		if err != nil {
			dukerror.Handle(dukerror.Normalize(err))
		}
	}()

	go func() {
		apiServer := api.NewServer(*apiAt, me)
		_ = apiServer.Start()
	}()

	go func() {
		me.StartSuspectedPeerChecker(time.Duration((*timerTicker) * int(time.Second)))
	}()

	select {}
}
