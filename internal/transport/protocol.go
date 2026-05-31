package transport

import (
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/baltej223/dukedb/internal/cluster"
	"github.com/baltej223/dukedb/scripts"
)

type MessageType int

const (
	PING MessageType = iota
	PONG

	PUT
	PUT_ACK
	PUT_REJ

	GET
	GET_RESPONSE

	GOSSIP

	JOIN
	JOIN_ACK
	JOIN_REJECT
)

func (m MessageType) String() string {
	switch m {
	case PING:
		return "PING"
	case PONG:
		return "PONG"
	case PUT:
		return "PUT"
	case PUT_ACK:
		return "PUT_ACK"
	case PUT_REJ:
		return "PUT_REJ"
	case GET:
		return "GET"
	case GET_RESPONSE:
		return "GET_RESPONSE"
	case GOSSIP:
		return "GOSSIP"
	case JOIN:
		return "JOIN"
	case JOIN_ACK:
		return "JOIN_ACK"
	case JOIN_REJECT:
		return "JOIN_REJECT"
	default:
		return "UNKNOWN"
	}
}

func ParseMessageType(s string) (MessageType, error) {
	switch strings.TrimSpace(s) {
	case "PING":
		return PING, nil
	case "PONG":
		return PONG, nil
	case "PUT":
		return PUT, nil
	case "PUT_ACK":
		return PUT_ACK, nil
	case "PUT_REJ":
		return PUT_REJ, nil
	case "GET":
		return GET, nil
	case "GET_RESPONSE":
		return GET_RESPONSE, nil
	case "GOSSIP":
		return GOSSIP, nil
	case "JOIN":
		return JOIN, nil
	case "JOIN_ACK":
		return JOIN_ACK, nil
	case "JOIN_REJECT":
		return JOIN_REJECT, nil
	default:
		return 0, fmt.Errorf("unknown message type: %s", s)
	}
}

type Peer = cluster.Peer

type Message struct {
	Type      MessageType
	RequestID string
	Headers   map[string]string
}

type ParsedMessage struct {
	Type      MessageType
	RequestID string

	NodeID string

	Key   string
	Value []byte

	Found   bool
	Success bool

	Reason string

	Peers []Peer
}

func (PM ParsedMessage) String() string {
	result := "TYPE " + PM.Type.String() + "\n"
	result += "REQUEST_ID " + PM.RequestID + "\n"

	return result
}

func EncodeValue(v []byte) string {
	return base64.StdEncoding.EncodeToString(v)
}

func DecodeValue(v string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(v)
}

// Serialize converts the Message struct to string
func Serialize(msg Message) string {
	var b strings.Builder

	b.WriteString(msg.Type.String())
	b.WriteByte('\n')

	b.WriteString("REQUEST_ID ")
	b.WriteString(msg.RequestID)
	b.WriteByte('\n')

	keys := make([]string, 0, len(msg.Headers))

	for k := range msg.Headers {
		if k == "REQUEST_ID" {
			continue
		}
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		b.WriteString(k)
		b.WriteByte(' ')
		b.WriteString(msg.Headers[k])
		b.WriteByte('\n')
	}

	return b.String()
}

// Parse converts the Message struct to string
func Parse(raw string) (ParsedMessage, error) {
	raw = strings.TrimSpace(raw)

	if raw == "" {
		return ParsedMessage{}, errors.New("empty message")
	}

	lines := strings.Split(raw, "\n")

	msgType, err := ParseMessageType(lines[0])
	if err != nil {
		return ParsedMessage{}, err
	}

	headers := make(map[string]string)

	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)

		if len(parts) != 2 {
			return ParsedMessage{}, fmt.Errorf("invalid line: %s", line)
		}

		headers[parts[0]] = parts[1]
	}

	msg := ParsedMessage{
		Type:      msgType,
		RequestID: headers["REQUEST_ID"],
	}

	switch msgType {

	case PING:
		msg.NodeID = headers["NODE_ID"]

	case PONG:
		msg.NodeID = headers["NODE_ID"]

	case JOIN:
		msg.NodeID = headers["NODE_ID"]

	case PUT:
		msg.Key = headers["KEY"]

		if encoded := headers["VALUE_BASE64"]; encoded != "" {
			msg.Value, err = DecodeValue(encoded)
			if err != nil {
				return ParsedMessage{}, err
			}
		}

	case GET:
		msg.Key = headers["KEY"]

	case GET_RESPONSE:
		msg.Found = headers["FOUND"] == "true"

		if encoded := headers["VALUE_BASE64"]; encoded != "" {
			msg.Value, err = DecodeValue(encoded)
			if err != nil {
				return ParsedMessage{}, err
			}
		}

	case PUT_ACK:
		msg.Success = headers["SUCCESS"] == "true"

	case PUT_REJ:
		msg.Success = headers["SUCCESS"] == "true"

	case JOIN_REJECT:
		msg.Reason = headers["REASON"]

	case JOIN_ACK:

		countStr := headers["PEER_COUNT"]

		count, err := strconv.Atoi(countStr)
		if err != nil {
			return ParsedMessage{}, err
		}

		for i := 0; i < count; i++ {

			nodeID := headers[fmt.Sprintf(
				"PEER_%d_NODE_ID",
				i,
			)]

			addr := headers[fmt.Sprintf(
				"PEER_%d_ADDR",
				i,
			)]

			msg.Peers = append(
				msg.Peers,
				Peer{
					NodeID: nodeID,
					Addr:   addr,
				},
			)
		}
	}

	return msg, nil
}

func createRequestID() (string, error) {
	id, err := scripts.RandomString(20)
	if err != nil {
		return "", err
	}

	if len(id) != 20 {
		return "", fmt.Errorf(
			"invalid request id length: got=%d want=20",
			len(id),
		)
	}

	return id, nil
}

// ----------------------------------------------------
// REQUEST BUILDERS
// ----------------------------------------------------

func CreatePingMessage(
	nodeID string,
) (Message, error) {
	requestID, err := createRequestID()
	if err != nil {
		return Message{}, err
	}

	return Message{
		Type:      PING,
		RequestID: requestID,
		Headers: map[string]string{
			"NODE_ID": nodeID,
		},
	}, nil
}

func CreatePutMessage(
	key string,
	value []byte,
) (Message, error) {
	requestID, err := createRequestID()
	if err != nil {
		return Message{}, err
	}

	return Message{
		Type:      PUT,
		RequestID: requestID,
		Headers: map[string]string{
			"KEY":          key,
			"VALUE_BASE64": EncodeValue(value),
		},
	}, nil
}

func CreateGetMessage(
	key string,
) (Message, error) {
	requestID, err := createRequestID()
	if err != nil {
		return Message{}, err
	}

	return Message{
		Type:      GET,
		RequestID: requestID,
		Headers: map[string]string{
			"KEY": key,
		},
	}, nil
}

func CreateJoinMessage(
	nodeID string,
	addr string,
) (Message, error) {
	requestID, err := createRequestID()
	if err != nil {
		return Message{}, err
	}

	return Message{
		Type:      JOIN,
		RequestID: requestID,
		Headers: map[string]string{
			"NODE_ID": nodeID,
			"ADDR":    addr,
		},
	}, nil
}

// ----------------------------------------------------
// RESPONSE BUILDERS
// ----------------------------------------------------

func CreatePongMessage(
	requestID string,
	nodeID string,
) Message {
	return Message{
		Type:      PONG,
		RequestID: requestID,
		Headers: map[string]string{
			"NODE_ID": nodeID,
		},
	}
}

func CreatePutACKMessage(
	requestID string,
) Message {
	return Message{
		Type:      PUT_ACK,
		RequestID: requestID,
		Headers: map[string]string{
			"SUCCESS": "true",
		},
	}
}

func CreatePutREJMessage(
	requestID string,
) Message {
	return Message{
		Type:      PUT_REJ,
		RequestID: requestID,
		Headers: map[string]string{
			"SUCCESS": "false",
		},
	}
}

func CreateGetResponseMessage(
	requestID string,
	found bool,
	value []byte,
) Message {
	return Message{
		Type:      GET_RESPONSE,
		RequestID: requestID,
		Headers: map[string]string{
			"FOUND":        strconv.FormatBool(found),
			"VALUE_BASE64": EncodeValue(value),
		},
	}
}

func CreateJoinACKMessage(
	requestID string,
	peers []Peer,
) Message {
	headers := map[string]string{
		"PEER_COUNT": strconv.Itoa(len(peers)),
	}

	for i, peer := range peers {

		headers[fmt.Sprintf(
			"PEER_%d_NODE_ID",
			i,
		)] = peer.NodeID

		headers[fmt.Sprintf(
			"PEER_%d_ADDR",
			i,
		)] = peer.Addr
	}

	return Message{
		Type:      JOIN_ACK,
		RequestID: requestID,
		Headers:   headers,
	}
}

func CreateJoinREJECTMessage(
	requestID string,
	reason string,
) Message {
	return Message{
		Type:      JOIN_REJECT,
		RequestID: requestID,
		Headers: map[string]string{
			"REASON": reason,
		},
	}
}
