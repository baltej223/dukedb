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
	Type      MessageType
	RequestID string
	headers   string
	Payload   string
}

type MessageResponse struct {
	Type            string
	RequestID       string
	ResponsePayload string
}

type MessageType int

const (
	PING MessageType = iota
	PONG
	PUT
	PUT_ACK
	GET
	GET_RESPONSE
	GOSSIP
	JOIN_ACK
	JOIN
	JOIN_REJECT
)

type ParsedMessage struct {
	Type   MessageType
	NodeID string

	RequestID string

	Key   string
	Value string

	Success bool
	Found   bool
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
	headers := "PING FROM \nREQUEST_ID " + id + "\n"
	return Message{PING, id, headers, ""}, nil
}

func CreatePongMessage() (Message, error) {
	id, err := scripts.RandomString(20)
	if err != nil {
		return Message{}, err
	}
	headers := "PONG FROM REQUEST_ID " + id + "\n"
	return Message{PONG, id, headers, ""}, nil
}

func CreateGetMessage(key string) (Message, error) {
	id, err := scripts.RandomString(20)
	if err != nil {
		return Message{}, err
	}
	headers := "GET\nREQUEST_ID " + id + "\nKEY " + key + "\n"
	return Message{GET, id, headers, ""}, nil
}

func CreateGetResponseMessage(key string, value string) (Message, error) {
	id, err := scripts.RandomString(20)
	if err != nil {
		return Message{}, err
	}
	headers := "GET RESPONSE\nREQUEST_IID " + id + "\nKEY " + key + "\n"
	return Message{GET_RESPONSE, id, headers, "VALUE " + value}, nil
}

func CreateGossipMessage(payload string) (Message, error) {
	id, err := scripts.RandomString(20)
	if err != nil {
		return Message{}, err
	}
	headers := "GOSSIP\nREQUEST_ID " + id + "\n"
	return Message{GOSSIP, id, headers, payload}, nil
}

func CreateJoinMessage(address string, value string) (Message, error) {
	nodeId, err := scripts.RandomAplhanumericString(10) // Node ids 10 chars alphanumeric
	if err != nil {
		return Message{}, err
	}
	id, err := scripts.RandomString(20)
	if err != nil {
		return Message{}, err
	}
	// ADDR 10.0.0.5:9000
	headers := "JOIN\nID " + nodeId + "\n REQUEST_ID " + id + "\nADDR" + address + "\n"
	return Message{JOIN, id, headers, "VALUE " + value}, nil
}

func CreateJoinACKMessage(address string) (Message, error) {
	nodeId, err := scripts.RandomAplhanumericString(10) // Node ids 10 chars alphanumeric
	if err != nil {
		return Message{}, err
	}
	id, err := scripts.RandomString(20)
	if err != nil {
		return Message{}, err
	}
	// ADDR 10.0.0.5:9000
	headers := "JOIN\nID " + nodeId + "\n REQUEST_ID " + id + "\nADDR" + address + "\n"
	return Message{JOIN_ACK, id, headers, ""}, nil
}

func CreateJoinREJMessage(address string, value string) (Message, error) {
	nodeId, err := scripts.RandomAplhanumericString(10) // Node ids 10 chars alphanumeric
	if err != nil {
		return Message{}, err
	}
	id, err := scripts.RandomString(20)
	if err != nil {
		return Message{}, err
	}
	// ADDR 10.0.0.5:9000
	headers := "JOIN\nID " + nodeId + "\n REQUEST_ID " + id + "\nADDR" + address + "\n"
	return Message{JOIN_REJECT, id, headers, "VALUE " + value}, nil
}

// func ParseMessage(raw string) (*Message, error) {
// 	lines := strings.Split(strings.TrimSpace(raw), "\n")
//
// 	if len(lines) == 0 {
// 		return nil, fmt.Errorf("empty message")
// 	}
//
// 	msg := &Message{}
//
// 	switch lines[0] {
//
// 	case "PUT":
// 		msg.Type = PUT
//
// 	case "PUT_ACK":
// 		msg.Type = PUT_ACK
//
// 	case "GET":
// 		msg.Type = GET
//
// 	case "GET RESPONSE":
// 		msg.Type = GET_RESPONSE
//
// 	default:
// 		if strings.HasPrefix(lines[0], "PING FROM ") {
// 			msg.Type = PING
// 			msg.NodeID = strings.TrimPrefix(lines[0], "PING FROM ")
// 			return msg, nil
// 		}
//
// 		if strings.HasPrefix(lines[0], "PONG FROM ") {
// 			msg.Type = PONG
// 			msg.NodeID = strings.TrimPrefix(lines[0], "PONG FROM ")
// 			return msg, nil
// 		}
//
// 		return nil, fmt.Errorf("unknown message type")
// 	}
//
// 	for _, line := range lines[1:] {
// 		parts := strings.SplitN(line, " ", 2)
//
// 		if len(parts) != 2 {
// 			continue
// 		}
//
// 		key := parts[0]
// 		value := parts[1]
//
// 		switch key {
//
// 		case "REQUEST_ID":
// 			msg.RequestID = value
//
// 		case "KEY":
// 			msg.Key = value
//
// 		case "VALUE":
// 			msg.Value = value
//
// 		case "SUCCESS":
// 			msg.Success = value == "true"
//
// 		case "FOUND":
// 			msg.Found = value == "true"
// 		}
// 	}
//
// 	return msg, nil
// } // SendMessage  the message and wraps the output int MessageResponse
//

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
* Send  It initiates a connections add length to it,
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
