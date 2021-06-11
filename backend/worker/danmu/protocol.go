package danmu

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	headerLength     int    = 16
	OpHeartbeat      uint32 = 2
	OpHeartbeatReply uint32 = 3
	OpSendSMSReply   uint32 = 5
	OpAuth           uint32 = 7
	OpAuthReply      uint32 = 8
	VerJSON          uint16 = 0
	VerZLIB          uint16 = 2
	seqenceID        uint32 = 0
)

/**
 * │--------│--------│--------│--------│
 * ┌───────────────────────────────────┐
 * │           Packet Length           │
 * ├─────────────────┬─────────────────┤
 * │  Header Length  │     Version     │
 * ├─────────────────┴─────────────────┤
 * │              Operation            │
 * ├───────────────────────────────────┤
 * │             Sequence ID           │
 * ├───────────────────────────────────┤
 * │                                   │
 * │                Body               │
 * │                                   │
 * └───────────────────────────────────┘
 */
type Message struct {
	packageLength uint32
	headerLength  uint16
	version       uint16
	Operation     uint32
	sequenceID    uint32
	Body          []byte
}

func decode(buffer []byte) []*Message {
	offset := 0
	messages := []*Message{}
	for offset < len(buffer) {
		message := Message{
			packageLength: binary.BigEndian.Uint32(buffer[offset+0:]),
			headerLength:  binary.BigEndian.Uint16(buffer[offset+4:]),
			version:       binary.BigEndian.Uint16(buffer[offset+6:]),
			Operation:     binary.BigEndian.Uint32(buffer[offset+8:]),
			sequenceID:    binary.BigEndian.Uint32(buffer[offset+12:]),
		}
		bodyBuffer := buffer[offset+int(message.headerLength) : offset+int(message.packageLength)]
		switch message.Operation {
		case OpAuthReply:
			// do nothing
		case OpHeartbeatReply:
			message.Body = bodyBuffer
		default:
			if message.version == VerJSON {
				message.Body = bodyBuffer
			}
			if message.version == VerZLIB {
				r, _ := zlib.NewReader(bytes.NewReader(bodyBuffer))
				b, _ := ioutil.ReadAll(r)
				messages = append(messages, decode(b)...)
			}
		}
		offset += int(message.packageLength)
		if message.version != VerZLIB {
			messages = append(messages, &message)
		}
	}
	return messages
}

func encode(op uint32, payload string) (buffer []byte) {
	packageLength := headerLength + len(payload)
	buff := make([]byte, packageLength)

	// package length
	binary.BigEndian.PutUint32(buff, uint32(packageLength))
	// header length
	binary.BigEndian.PutUint16(buff[4:], uint16(headerLength))
	// version
	binary.BigEndian.PutUint16(buff[6:], VerJSON)
	// operation
	binary.BigEndian.PutUint32(buff[8:], op)
	// sequence id
	binary.BigEndian.PutUint32(buff[12:], seqenceID)
	// body
	copy(buff[headerLength:], payload)

	return buff
}

func joinRoom(roomID int64, uid int) []byte {
	payload := fmt.Sprintf(`{"uid":%d,"roomid":%d,"protover":2,"platform":"3rd_party"}`, uid, roomID)
	return encode(OpAuth, payload)
}

func heartbeat() []byte {
	return encode(OpHeartbeat, "")
}

func (m *Message) String() string {
	switch m.Operation {
	case OpSendSMSReply:
		b := bytes.Buffer{}
		json.Indent(&b, m.Body, "", "\t")
		return fmt.Sprintf("SMS_REPLY\n%s\n", b.String())
	case OpAuthReply:
		return "Auth\n"
	case OpHeartbeatReply:
		return fmt.Sprintf("HEARTBEAT\nonline: %d\n", binary.BigEndian.Uint32(m.Body))
	default:
		return "unidentified message type"
	}
}
