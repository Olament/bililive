package worker

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const (
	heartbeatInv   = time.Second * 10
	reconnectDelay = time.Second * 3
	timeout        = time.Second * 30
)

func randomID() int {
	return 1e15 + int(rand.Float32()*2e15)
}

func connect(ctx context.Context, roomID int64, out chan *message) error {
	u := url.URL{
		Scheme: "wss",
		Host:   "broadcastlv.chat.bilibili.com:2245",
		Path:   "sub",
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	// send join room request
	conn.SetWriteDeadline(time.Now().Add(timeout))
	err = conn.WriteMessage(websocket.BinaryMessage, joinRoom(roomID, randomID()))
	if err != nil {
		return err
	}

	// handle inbound message and periodically do heartbeat
	heartbeat := heartbeat()
	ticker := time.NewTicker(heartbeatInv)
	for {
		select {
		case <-ctx.Done():
			close(out)
			return nil
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(timeout))
			conn.WriteMessage(websocket.BinaryMessage, heartbeat)
		default:
			conn.SetReadDeadline(time.Now().Add(timeout))
			_, b, err := conn.ReadMessage()
			if err != nil {
				return err
			}
			for _, msg := range decode(b) {
				out <- msg
			}
		}
	}
}

// Connect is a blocking function that reads the message from broadcast with roomID and
// then push it to the out channel
func Connect(ctx context.Context, roomID int64, out chan *message) {
	err := connect(ctx, roomID, out)
	if err != nil {
		fmt.Printf("worker/websocket: [%d] %v\n", roomID, err)
	}
	select {
	case <-ctx.Done():
		return
	default:
		time.AfterFunc(reconnectDelay+time.Duration(rand.Int31n(100)), func() {
			fmt.Printf("worker/websocket: [%d] reconnect...\n", roomID)
			Connect(ctx, roomID, out)
		})
	}
}
