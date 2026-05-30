package transport

import (
	"bufio"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/baltej223/dukedb/internal/cluster"
)

func SendMessage(n cluster.Peer, m Message) (ParsedMessage, error) {
	messageString := Serialize(m)
	response, err := Send(n, messageString)
	if err != nil {
		return ParsedMessage{},
			err
	}
	parsedResponse, err := Parse(response)
	if err != nil {
		return ParsedMessage{},
			err
	}
	return parsedResponse, nil
}

/*
* Send function sends the provided string to the node
* */

func Send(n cluster.Peer, payload string) (string, error) {
	length := strconv.Itoa(len(payload))

	// Initiate the connection
	conn, err := net.Dial("tcp", n.Addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	packet := []byte(length + "\n" + payload + "\n\n")

	_, err = conn.Write(packet)
	if err != nil {
		return "", err
	}

	response, error := readMessage(conn)
	if error != nil {
		return "", error
	}

	return response, nil
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
