package worker

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
	opHeartbeat      uint32 = 2
	opHeartbeatReply uint32 = 3
	opSendSMSReply   uint32 = 5
	opAuth           uint32 = 7
	opAuthReply      uint32 = 8
	verJSON          uint16 = 0
	verZLIB          uint16 = 2
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
type message struct {
	packageLength uint32
	headerLength  uint16
	version       uint16
	operation     uint32
	sequenceID    uint32
	body          []byte
}

func decode(buffer []byte) []*message {
	offset := 0
	messages := []*message{}
	for offset < len(buffer) {
		message := message{
			packageLength: binary.BigEndian.Uint32(buffer[offset+0:]),
			headerLength:  binary.BigEndian.Uint16(buffer[offset+4:]),
			version:       binary.BigEndian.Uint16(buffer[offset+6:]),
			operation:     binary.BigEndian.Uint32(buffer[offset+8:]),
			sequenceID:    binary.BigEndian.Uint32(buffer[offset+12:]),
		}
		bodyBuffer := buffer[offset+int(message.headerLength) : offset+int(message.packageLength)]
		switch message.operation {
		case opAuthReply:
			// do nothing
		case opHeartbeatReply:
			message.body = bodyBuffer
		default:
			if message.version == verJSON {
				message.body = bodyBuffer
			}
			if message.version == verZLIB {
				r, _ := zlib.NewReader(bytes.NewReader(bodyBuffer))
				b, _ := ioutil.ReadAll(r)
				messages = append(messages, decode(b)...)
			}
		}
		offset += int(message.packageLength)
		if message.version != verZLIB {
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
	binary.BigEndian.PutUint16(buff[6:], verJSON)
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
	return encode(opAuth, payload)
}

func heartbeat() []byte {
	return encode(opHeartbeat, "")
}

func (m *message) String() string {
	switch m.operation {
	case opSendSMSReply:
		b := bytes.Buffer{}
		json.Indent(&b, m.body, "", "\t")
		return fmt.Sprintf("SMS_REPLY\n%s\n", b.String())
	case opAuthReply:
		return "Auth\n"
	case opHeartbeatReply:
		return fmt.Sprintf("HEARTBEAT\nonline: %d\n", binary.BigEndian.Uint32(m.body))
	default:
		return "unidentified message type"
	}
}