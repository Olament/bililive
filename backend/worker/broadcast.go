package worker

import (
	"bililive/worker/common"
	"bililive/worker/danmu"
	"context"
	"encoding/binary"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/tidwall/gjson"
)

type Broadcast struct {
	Roomid                 int64          `json:"roomid"`
	UID                    int64          `json:"uid"`
	Uname                  *common.String `json:"uname"`
	Popularity             uint32         `json:"popularity"`
	MaxPopularity          uint32         `json:"maxPopularity"`
	Title                  *common.String `json:"title"`
	Usercover              *common.String `json:"usercover"`
	Keyframe               *common.String `json:"keyframe"`
	Livetime               time.Time      `json:"livetime"`
	Endtime                time.Time      `json:"-"`
	Participantduring10Min int64          `json:"participantDuring10Min"`
	GoldCoin               uint64         `json:"goldCoin"`
	SilverCoin             uint64         `json:"silverCoin"`

	cancel context.CancelFunc
	isStop uint32
	setTTL *common.Set
}

func (b *Broadcast) start() {
	out := make(chan *danmu.Message, 10)
	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel
	b.setTTL = common.NewSet(time.Minute * 10)

	go danmu.Connect(ctx, b.Roomid, out)
	for msg := range out {
		b.parseMessage(msg)
	}
}

func (b *Broadcast) stop() {
	if ok := atomic.CompareAndSwapUint32(&b.isStop, 0, 1); ok {
		b.cancel()
		b.Endtime = time.Now()
	}
}

func (b *Broadcast) parseMessage(msg *danmu.Message) {
	switch msg.Operation {
	case danmu.OpHeartbeatReply:
		popularity := binary.BigEndian.Uint32(msg.Body)
		atomic.StoreUint32(&b.Popularity, popularity)
		if popularity == 1 {
			b.stop()
		}
		if popularity > b.MaxPopularity {
			atomic.StoreUint32(&b.MaxPopularity, popularity)
		}
	case danmu.OpSendSMSReply:
		switch gjson.GetBytes(msg.Body, "cmd").String() {
		case "COMBO_SEND", "SEND_GIFT":
			res := gjson.GetManyBytes(msg.Body, "data.coin_type", "data.total_coin")
			if res[0].String() == "silver" {
				atomic.AddUint64(&b.SilverCoin, res[1].Uint())
			} else {
				atomic.AddUint64(&b.GoldCoin, res[1].Uint())
			}
		case "GUARD_BUY", "SUPER_CHAT_MESSAGE":
			res := gjson.GetBytes(msg.Body, "data.price")
			atomic.AddUint64(&b.GoldCoin, res.Uint())
		case "DANMU_MSG":
			uid := gjson.GetBytes(msg.Body, "info.2.0").Int()
			b.setTTL.Add(uid)
			atomic.StoreInt64(&b.Participantduring10Min, b.setTTL.Len())
		case "PREPARING":
			b.stop()
		}
	case danmu.OpAuthReply:
	default:
		fmt.Println("worker/broadcast: unidentified message type")
	}
}

// return a Copy of broadcast atomically
// only the public field will be copied and it should only be used for json Marshal
func (b *Broadcast) Copy() *Broadcast {
	broadcast := &Broadcast{
		Roomid:                 atomic.LoadInt64(&b.Roomid),
		UID:                    atomic.LoadInt64(&b.UID),
		Uname:                  b.Uname,
		Popularity:             atomic.LoadUint32(&b.Popularity),
		MaxPopularity:          atomic.LoadUint32(&b.MaxPopularity),
		Title:                  b.Title,
		Usercover:              b.Usercover,
		Keyframe:               b.Keyframe,
		Livetime:               b.Livetime,
		Endtime:                b.Endtime,
		Participantduring10Min: atomic.LoadInt64(&b.Participantduring10Min),
		GoldCoin:               atomic.LoadUint64(&b.GoldCoin),
		SilverCoin:             atomic.LoadUint64(&b.SilverCoin),
	}
	return broadcast
}
