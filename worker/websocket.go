package worker

import (
	"math/rand"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const (
	heartbeatInv   = time.Second * 10
	reconnectDelay = time.Second * 3
)

func randomID() int {
	return 1e15 + int(rand.Float32()*2e15)
}

func connect(roomID int) error {
	u := url.URL{
		Scheme: "wss",
		Host:   "broadcastlv.chat.bilibili.com:2245",
		Path:   "sub",
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	defer conn.Close()
	if err != nil {
		return err
	}

	// send join room request
	err = conn.WriteMessage(websocket.BinaryMessage, joinRoom(roomID, randomID()))
	if err != nil {
		return err
	}

	// handle heartbeat
	heartbeat := heartbeat()
	go func() {
		for range time.Tick(heartbeatInv) {
			conn.WriteMessage(websocket.BinaryMessage, heartbeat)
		}
	}()

	// handle inbound message
	for {
		_, b, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		for _, msg := range decode(b) {
			parseMessage(msg)
		}
	}
}

func SimpleConnect() {
	connect(22333522)
}
