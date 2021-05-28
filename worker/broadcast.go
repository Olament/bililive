package worker

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type broadcast struct {
	Roomid                 int
	Shortid                int
	UID                    int
	Uname                  string
	Popularity             int
	Title                  string
	Usercover              string
	Keyframe               string
	Livetime               string
	Participantduring10Min int

	cancel context.CancelFunc
}

func (b *broadcast) start() {
	out := make(chan *message, 10)
	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel
	go Connect(ctx, b.Roomid, out)
	for msg := range out {
		b.parseMessage(msg)
	}
}

func (b *broadcast) stop() {
	b.cancel()
}

func (b *broadcast) parseMessage(msg *message) {
	switch msg.operation {
	case opHeartbeatReply:
		fmt.Printf("HEARTBEAT\nonline: %d\n\n", binary.BigEndian.Uint32(msg.body))
	case opSendSMSReply:
		buffer := bytes.Buffer{}
		json.Indent(&buffer, msg.body, "", "\t")
		fmt.Printf("SMS_REPLY\ndata: %s\n\n", buffer.String())
	case opAuthReply:
		fmt.Printf("AUTH\n\n")
	default:
		fmt.Println("worker/protocol: unidentified message type")
	}
}
