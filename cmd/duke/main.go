package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/storing"
	"github.com/baltej223/dukedb/internal/transport"
)

func main() {
	// Flags handling
	selfAddress := flag.String("selfAddr", "localhost:8000", "Address of the current node, Example localhost:8000")
	peerAddress := flag.String("peerAddr", "localhost:8001", "Address of peer node, Example: localhost:8001")
	peerNodeID := flag.String("peerNodeID", "b", "Peer node ID, Example: b")
	// FLAGS END
	flag.Parse()

	hostname := *selfAddress
	GloabalHOSTNAME := hostname
	server := transport.NewServer(hostname)

	// Set up internal KV
	storing.InitialiseKVI()
	neighbours := []cluster.Peer{{*peerNodeID, *peerAddress}}
	neighboursBytes, err := json.Marshal(neighbours)
	if err != nil {
		panic(err) // or handle the error properly
	}
	storing.PutI("neighbours", neighboursBytes)

	log.Println("Starting duke node on " + "8000")

	go func() {
		err := server.Start(transport.HandleMessage)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done")
	}()

	time.Sleep(5 * time.Second)

	message, _ := transport.CreatePingMessage("a")
	p := neighbours[0]

	_ = transport.SendMessage(p, message, GloabalHOSTNAME)
	select {}
}
