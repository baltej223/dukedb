package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/node"
	"github.com/baltej223/dukedb/internal/storing"
	"github.com/baltej223/dukedb/internal/transport"
	"github.com/baltej223/dukedb/scripts"
)

func main() {
	// Flags handling
	selfAddress := flag.String("self-addr", "localhost:8000", "Address of the current node, Example localhost:8000")
	selfNodeID := flag.String("self-node-id", "-", "Self node ID, Example: b")

	isSeedNode := flag.Bool("seed-node", false, "Is this a seed node.")

	peerAddress := flag.String("peer-addr", "", "Address of peer node, Example: localhost:8001")
	peerNodeID := flag.String("peer-node-id", "", "Peer node ID, Example: b")
	delay := flag.Int("delay", 5, "[Debug]: Initial Delay Before sending first request")

	forPutGet := flag.Bool("yaay", false, "pingu")
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

	// Start gossip loop now

	go func() {
		err := me.StartGossipLoop(true)
		if err != nil {
			log.Printf("gossip failed: %v", err)
		}
	}()

	go func() {
		if *forPutGet {
			for {
				log.Printf("[TEST] starting put/get test")

				time.Sleep(15 * time.Second)

				log.Printf("[TEST] woke up after sleep")

				peers := me.Cluster.GetPeers()
				log.Printf("[TEST] cluster has %d peers", len(peers))

				randomNode, err := scripts.ChooseRandomElements(
					peers,
					1,
				)
				if err != nil {
					log.Printf(
						"[TEST] failed to choose random peer: %v",
						err,
					)
					return
				}

				target := randomNode[0]

				log.Printf(
					"[TEST] selected target node=%s addr=%s",
					target.NodeID,
					target.Addr,
				)

				putMessage, err := transport.CreatePutMessage(
					"Name",
					[]byte("Baltej"),
					me.ID,
				)
				fmt.Println(transport.Serialize(putMessage))
				if err != nil {
					log.Printf(
						"[TEST] failed to create PUT: %v",
						err,
					)
					return
				}

				log.Printf(
					"[TEST] sending PUT request_id=%s",
					putMessage.RequestID,
				)

				putResponse, err := me.SendRequestAndWait(
					target,
					putMessage,
					10*time.Second,
				)
				if err != nil {
					log.Printf(
						"[TEST] PUT failed: %v",
						err,
					)
					return
				}

				log.Printf(
					"[TEST] PUT response received type=%s request_id=%s",
					putResponse.Type.String(),
					putResponse.RequestID,
				)

				getReq, err := transport.CreateGetMessage(
					"Name",
					me.ID,
				)
				if err != nil {
					log.Printf(
						"[TEST] failed to create GET: %v",
						err,
					)
					return
				}

				log.Printf(
					"[TEST] sending GET request_id=%s",
					getReq.RequestID,
				)

				response, err := me.SendRequestAndWait(
					target,
					getReq,
					20*time.Second,
				)
				if err != nil {
					log.Printf(
						"[TEST] GET failed: %v",
						err,
					)
					return
				}

				log.Printf(
					"[TEST] GET response received request_id=%s value=%s",
					response.RequestID,
					string(response.Value),
				)
			}
		}
	}()

	select {}
}
