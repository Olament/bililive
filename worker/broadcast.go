package worker

import (
	"context"
	"encoding/binary"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/tidwall/gjson"
)

type Broadcast struct {
	Roomid                 int64     `json:"roomid"`
	UID                    int64     `json:"uid"`
	Uname                  string    `json:"uname"`
	Popularity             uint32    `json:"popularity"`
	MaxPopularity          uint32    `json:"max_popularity"`
	Title                  string    `json:"title"`
	Usercover              string    `json:"usercover"`
	Keyframe               string    `json:"keyframe"`
	Livetime               time.Time `json:"livetime"`
	Endtime                time.Time `json:"endtime"`
	Participantduring10Min int       `json:"participantduring_10_min"`
	GoldCoin               uint64    `json:"gold_coin"`
	SilverCoin             uint64    `json:"silver_coin"`

	cancel context.CancelFunc
	isStop uint32
}

func (b *Broadcast) start() {
	out := make(chan *message, 10)
	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel
	go Connect(ctx, b.Roomid, out)
	for msg := range out {
		b.parseMessage(msg)
	}
}

func (b *Broadcast) stop() {
	if ok := atomic.CompareAndSwapUint32(&b.isStop, 0, 1); ok {
		b.cancel()
		b.Endtime = time.Now()
		fmt.Printf("%+v\n", b)
	}
}

func (b *Broadcast) parseMessage(msg *message) {
	switch msg.operation {
	case opHeartbeatReply:
		popularity := binary.BigEndian.Uint32(msg.body)
		atomic.StoreUint32(&b.Popularity, popularity)
		if popularity == 1 {
			b.stop()
		}
		if popularity > b.MaxPopularity {
			atomic.StoreUint32(&b.MaxPopularity, popularity)
		}
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

func (b *Broadcast) Copy() *Broadcast {
	broadcast := *b
	return &broadcast
}
