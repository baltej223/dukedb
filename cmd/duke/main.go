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
		if !*forPutGet {
			return
		}

		testCounter := 0

		for {

			time.Sleep(15 * time.Second)

			peers := me.Cluster.GetPeers()

			randomNode, err := scripts.ChooseRandomElements(
				peers,
				1,
			)
			if err != nil {
				log.Printf(
					"[TEST] failed to choose random peer: %v",
					err,
				)
				continue
			}

			target := randomNode[0]

			key := fmt.Sprintf(
				"test-key-%d",
				testCounter,
			)

			expectedValue := fmt.Sprintf(
				"test-value-%d",
				testCounter,
			)

			log.Printf(
				"[TEST %d] PUT key=%s value=%s target=%s",
				testCounter,
				key,
				expectedValue,
				target.NodeID,
			)

			putMessage, err := transport.CreatePutMessage(
				key,
				[]byte(expectedValue),
				me.ID,
				me.MembershipVersion,
			)
			if err != nil {
				log.Printf(
					"[TEST %d] failed creating PUT: %v",
					testCounter,
					err,
				)
				continue
			}

			putResp, err := me.SendRequestAndWait(
				target,
				putMessage,
				10*time.Second,
			)
			if err != nil {
				log.Printf(
					"[TEST %d] PUT failed: %v",
					testCounter,
					err,
				)
				continue
			}

			log.Printf(
				"[TEST %d] PUT response=%s",
				testCounter,
				putResp.Type.String(),
			)

			getReq, err := transport.CreateGetMessage(
				key,
				me.ID,
				me.MembershipVersion,
			)
			if err != nil {
				log.Printf(
					"[TEST %d] failed creating GET: %v",
					testCounter,
					err,
				)
				continue
			}

			getResp, err := me.SendRequestAndWait(
				target,
				getReq,
				20*time.Second,
			)
			if err != nil {
				log.Printf(
					"[TEST %d] GET failed: %v",
					testCounter,
					err,
				)
				continue
			}

			actualValue := string(getResp.Value)

			if actualValue != expectedValue {
				log.Printf(
					"[TEST %d] MISMATCH expected=%q got=%q",
					testCounter,
					expectedValue,
					actualValue,
				)
			} else {
				log.Printf(
					"[TEST %d] SUCCESS key=%s value=%s",
					testCounter,
					key,
					actualValue,
				)
			}

			testCounter++
		}
	}()

	select {}
}
