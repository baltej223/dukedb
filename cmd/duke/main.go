package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/transport"
)

func handler(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()

	fmt.Println("new connection from", remoteAddr)

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("connection closed:", remoteAddr)
			return
		}

		fmt.Printf("received from %s: %s", remoteAddr, message)

		req, _ := transport.Parse(message)

		m := transport.CreatePongMessage(req.RequestID, req.NodeID)
		response := transport.Serialize(m)

		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("write error:", err)
			return
		}
	}
}

func main() {
	// a := node.Initialise("a", "localhost", "8080")
	// b := node.Initialise("b", "localhost", "8081")
	// network := make(map[string][]*node.Node)
	// network["a"] = []*node.Node{b} // b is connected to automatic
	// network["b"] = []*node.Node{a}
	//
	// me := a
	// // Eventually I will need to set it up for automatic ID choosing.

	hostname := ":8000"
	server := transport.NewServer(hostname)

	log.Println("Starting duke node on " + "8000")

	go func() {
		err := server.Start(handler)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done")
	}()

	time.Sleep(5 * time.Second)

	message, _ := transport.CreatePingMessage("a")
	p := cluster.Peer{"a", "localhost:8001"}

	response, _ := transport.SendMessage(p, message)
	fmt.Printf("%s", response.RequestID)
	select {}
}
