package main

import (
	"fmt"
	"log"
	"time"

	"github.com/baltej223/dukedb/internal/node"
	"github.com/baltej223/dukedb/internal/transport"
)

func main() {
	a := node.Initialise("a", "localhost", "8080")
	_ = a
	b := node.Initialise("b", "localhost", "8081")
	network := make(map[string][]*node.Node)
	network["a"] = []*node.Node{b} // b is connected to automatic
	network["b"] = []*node.Node{a}

	me := a
	// Eventually I will need to set it up for automatic ID choosing.

	hostname := me.GetFullHostname()
	server := transport.NewServer(hostname)

	log.Println("starting duke node on " + me.GetPort())

	go func() {
		err := server.Start(transport.HandleConnection)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done")
	}()

	time.Sleep(5 * time.Second)

	value, err := transport.Send(network[me.ID][0], "Hello from "+me.ID)
	if err != nil {
		panic(err)
	}

	fmt.Println(value)

	select {}
}
