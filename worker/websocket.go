package worker

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
)

func SimpleConnect() {
	u := url.URL{
		Scheme: "wss",
		Host: "broadcastlv.chat.bilibili.com:2245",
		Path: "sub",
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}

	roomID := 915905
	randomID := 2007177791539429
	conn.WriteMessage(websocket.BinaryMessage, joinRoom(roomID, randomID))
	go func() {
		conn.WriteMessage(websocket.BinaryMessage, heartbeat())
		for range time.Tick(time.Second * 10) {
			err := conn.WriteMessage(websocket.BinaryMessage, heartbeat())
			if err != nil {
				panic(err)
			}
			fmt.Println("heartbeat!")
		}
	}()

	for {
		_, b, err := conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		for _, msg := range decode(b) {
			fmt.Println(msg)
		}
	}
}
