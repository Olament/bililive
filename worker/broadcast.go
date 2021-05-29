package worker

import (
	"context"
	"encoding/binary"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/tidwall/gjson"
)

type broadcast struct {
	Roomid                 int64
	UID                    int64
	Uname                  string
	Popularity             uint32
	Title                  string
	Usercover              string
	Keyframe               string
	Livetime               time.Time
	Participantduring10Min int
	GoldCoin               uint64
	SilverCoin             uint64

	cancel context.CancelFunc
	isStop uint32
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
	if ok := atomic.CompareAndSwapUint32(&b.isStop, 0, 1); ok {
		b.cancel()
		fmt.Printf("%+v\n", b)
	}
}

func (b *broadcast) parseMessage(msg *message) {
	switch msg.operation {
	case opHeartbeatReply:
		popularity := binary.BigEndian.Uint32(msg.body)
		if popularity == 1 {
			b.stop()
		}
		atomic.StoreUint32(&b.Popularity, popularity)
	case opSendSMSReply:
		switch gjson.GetBytes(msg.body, "cmd").String() {
		case "COMBO_SEND", "SEND_GIFT":
			res := gjson.GetManyBytes(msg.body, "data.coin_type", "data.total_coin")
			if res[0].String() == "silver" {
				atomic.AddUint64(&b.SilverCoin, res[1].Uint())
			} else {
				atomic.AddUint64(&b.GoldCoin, res[1].Uint())
			}
		case "GUARD_BUY", "SUPER_CHAT_MESSAGE":
			res := gjson.GetBytes(msg.body, "data.price")
			atomic.AddUint64(&b.GoldCoin, res.Uint())
		case "DANMU_MSG":
		// TODO: handle DanMu
		case "PREPARING":
			b.stop()
		}
	case opAuthReply:
	default:
		fmt.Println("worker/broadcast: unidentified message type")
	}
}
