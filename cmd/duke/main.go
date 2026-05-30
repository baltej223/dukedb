package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/baltej223/dukedb/internal/storing"
	"github.com/baltej223/dukedb/internal/transport"
)

func handler(con net.Conn) {
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

		response := "ack: " + message

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
		err := server.Start(transport.HandleConnection)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done")
	}()

	time.Sleep(5 * time.Second)

	_ = storing.InitialiseKV()
	storing.Put("Name", []byte("Baltej"))

	get, _ := storing.Get("Name")
	fmt.Printf("Name: %s", get)

	select {}
}
