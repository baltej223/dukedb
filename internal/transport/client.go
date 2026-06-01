package transport

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/internal/storing"
)

func SendMessage(n cluster.Peer, m Message) error {
	hostname, _ := storing.GetI("hostname")
	messageString := Serialize(m)
	fmt.Printf("[%s] Sending message to %s\n", hostname, n.Addr)
	err := Send(n, messageString)
	if err != nil {
		return err
	}
	return nil
}

/*
* Send function sends the provided string to the node
* */

func Send(n cluster.Peer, payload string) error {
	length := strconv.Itoa(len(payload))

	// Initiate the connection
	conn, err := net.Dial("tcp", n.Addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	packet := []byte(length + "\n" + payload + "\n\n")

	_, err = conn.Write(packet)
	if err != nil {
		return err
	}
	return nil
}

/*
* Since the protocol defines the message as [length][message]
* So its necessary to send and receive the length of the message.
* Exactly n bytes are read.
 */
func readMessage(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)

	// Read length line
	lengthLine, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	lengthLine = strings.TrimSpace(lengthLine)

	length, err := strconv.Atoi(lengthLine)
	if err != nil {
		return "", err
	}

	payload := make([]byte, length)

	_, err = io.ReadFull(reader, payload)
	if err != nil {
		return "", err
	}

	return string(payload), nil
}
