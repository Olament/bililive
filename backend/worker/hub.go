package worker

import (
	"bililive/worker/common"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tidwall/gjson"
)

const (
	baseURL = `https://api.live.bilibili.com/room/v3/area/getRoomList`
)

type Hub struct {
	broadcasts sync.Map // roomID -> *broadcast
	hClient    *http.Client
	dClient    *mongo.Client
	ctx        context.Context
}

func (h *Hub) Init() {
	h.hClient = &http.Client{}
	h.broadcasts = sync.Map{}
	h.ctx = context.Background()

	// setup mongodb connection
	client, err := mongo.Connect(h.ctx, options.Client().ApplyURI(os.Getenv("DB_STR")))
	if err != nil {
		panic(err)
	}
	h.dClient = client

	h.update()
	go func() {
		for _ = range time.Tick(time.Minute * 1) {
			h.update()
		}
	}()
}

func (h *Hub) newBroadcast(res gjson.Result) *Broadcast {
	b := Broadcast{
		Roomid:        res.Get("roomid").Int(),
		UID:           res.Get("uid").Int(),
		Uname:         common.NewString(res.Get("uname").String()),
		Face:          common.NewString(res.Get("face").String()),
		Title:         common.NewString(res.Get("title").String()),
		Usercover:     common.NewString(res.Get("cover").String()),
		Keyframe:      common.NewString(res.Get("system_cover").String()),
		Livetime:      time.Now(),
		Popularity:    uint32(res.Get("online").Int()),
		MaxPopularity: uint32(res.Get("online").Int()),
		ctx:           h.ctx,
		collection:    h.dClient.Database("livevup").Collection("broadcast"),
	}
	go b.start()
	return &b
}

func (h *Hub) update() {
	list := []gjson.Result{}
	pageNum := 1
	for res := h.fetch(pageNum); len(res) > 0; res = h.fetch(pageNum) {
		list = append(list, res...)
		pageNum += 1
	}
	// add new broadcast to the map
	// update keyframe of the existing broadcast
	for _, res := range list {
		roomID := res.Get("roomid")
		if v, ok := h.broadcasts.Load(roomID); ok {
			v.(*Broadcast).Keyframe = common.NewString(res.Get("system_cover").String())
		} else {
			h.broadcasts.Store(roomID, h.newBroadcast(res))
		}
	}
	// removing stopped broadcast
	h.broadcasts.Range(func(key, value interface{}) bool {
		if atomic.LoadUint32(&value.(*Broadcast).isStop) == 1 {
			h.broadcasts.Delete(key)
		}
		return true
	})
}

func (h *Hub) fetch(page int) []gjson.Result {
	resp, err := h.hClient.Get(fmt.Sprintf("%s?parent_area_id=%d&page=%d&page_size=%d",
		baseURL, 9, page, 99))
	if err != nil {
		fmt.Printf("worker/hub: %v\n", err)
		return nil
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("worker/hub: %v\n", nil)
		return nil
	}
	return gjson.GetBytes(bs, "data.list").Array()
}
