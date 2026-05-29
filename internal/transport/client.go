package transport

import (
	"bufio"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/scripts"
)

type Message struct {
	Type      string
	RequestID string
	headers   string
	Payload   string
}

type MessageResponse struct {
	Type            string
	RequestID       string
	ResponsePayload string
}

func CreateMessage(Type string, Payload string) {
	allTypes := []string{"PING", "PONG", "JOIN", "JOIN_ACK", "PUT", "PUT_ACK", "GET", "GET_RESPONSE", "GOSSIP"}
	_ = allTypes
}

func CreatePingMessage() (Message, error) {
	id, err := scripts.RandomString(20)
	if err != nil {
		return Message{}, err
	}
	headers := "PING\nID " + id + "\n"
	return Message{"PING", id, headers, ""}, nil
}

func CreatePongMessage() (Message, error) {
	id, err := scripts.RandomString(20)
	if err != nil {
		return Message{}, err
	}
	headers := "PONG\nID " + id + "\n"
	return Message{"PONG", id, headers, ""}, nil
}

func CreateGetMessage(key string) (Message, error) {
	id, err := scripts.RandomString(20)
	if err != nil {
		return Message{}, err
	}
	headers := "GET\nID " + id + "\nKEY " + key + "\n"
	return Message{"GET", id, headers, ""}, nil
}

func CreateGetResponseMessage(key string, value string) (Message, error) {
	id, err := scripts.RandomString(20)
	if err != nil {
		return Message{}, err
	}
	headers := "GET RESPONSE\nID " + id + "\nKEY " + key + "\n"
	return Message{"GET RESPONSE", id, headers, "VALUE " + value}, nil
}

func CreateGossipMessage(payload string) (Message, error) {
	id, err := scripts.RandomString(20)
	if err != nil {
		return Message{}, err
	}
	headers := "GOSSIP\nID " + id + "\n"
	return Message{"GOSSIP", id, headers, payload}, nil
}

// SendMessage: the message and wraps the output int MessageResponse
func (m *Message) SendMessage(n *cluster.OtherNode) (MessageResponse, error) {
	stringMessage := m.headers + m.Payload

	result, err := Send(n, stringMessage)
	_ = result
	response := MessageResponse{m.Type, m.RequestID, result}

	if err != nil {
		return MessageResponse{}, nil
	}
	//
	// if m.Type != "GOSSIP" {
	// 	switch m.Type {
	// 	case "JOIN":
	// 		cluster.HandleJoin(result)
	// 	case "PUT":
	// 		cluster.HandlePut(result)
	// 	}
	// }

	return response, nil
}

/*
* Send : It initiates a connections add length to it,
* and sends a message to it, expects response.
* */
func Send(n *cluster.OtherNode, payload string) (string, error) {
	length := strconv.Itoa(len(payload))

	// Initiate the connection
	conn, err := net.Dial("tcp", n.Address)
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
