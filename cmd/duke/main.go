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
)

func main() {
	// Flags handling
	selfAddress := flag.String("selfAddr", "localhost:8000", "Address of the current node, Example localhost:8000")
	selfNodeID := flag.String("selfNodeID", "-", "Self node ID, Example: b")
	peerAddress := flag.String("peerAddr", "localhost:8001", "Address of peer node, Example: localhost:8001")
	peerNodeID := flag.String("peerNodeID", "b", "Peer node ID, Example: b")
	delay := flag.Int("delay", 5, "[Debug]: Initial Delay Before sending first request")
	// FLAGS END
	flag.Parse()

	hostname := *selfAddress

	// Build node here
	me := node.Initialise(*selfNodeID, *selfAddress)

	server := transport.NewServer(hostname)

	// Set up internal KV
	storing.InitialiseKVI()
	neighbours := []cluster.Peer{{*peerNodeID, *peerAddress}}
	storing.PutIJSON("neighbours", neighbours)

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

	go func() {
		time.Sleep(time.Duration(*delay) * time.Second)

		message, _ := transport.CreatePingMessage(me.ID)
		p := neighbours[0]

		response, _ := me.SendRequestAndWait(p, message, 10*time.Second)
		fmt.Println("Done " + response.Type.String())
	}()
	select {}
}
