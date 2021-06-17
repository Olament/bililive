package worker

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type Payload struct {
	Count int          `json:"count"`
	List  []*Broadcast `json:"list"`
}

type ShortBroadcast struct {
	Title            string    `json:"title" bson:"title"`
	MaxPopularity    uint32    `json:"maxPopularity" bson:"maxPopularity"`
	Livetime         time.Time `json:"livetime" bson:"livetime"`
	Endtime          time.Time `json:"endTime" bson:"endtime"`
	Participant      int64     `json:"participant" bson:"participant"`
	GoldCoin         uint64    `json:"goldCoin" bson:"goldCoin"`
	GoldUser         int64     `json:"goldUser" bson:"goldUser"`
	SilverCoin       uint64    `json:"silverCoin" bson:"silverCoin"`
	DanmuCount       uint64    `json:"danmuCount" bson:"danmuCount"`
	ParticipantTrend []int64   `json:"participantTrend" bson:"participantTrend"`
	GoldTrend        []uint64  `json:"goldTrend" bson:"goldTrend"`
	DanmuTrend       []uint64  `json:"danmuTrend" bson:"danmuTrend"`
}

func (h *Hub) Online() gin.HandlerFunc {
	f := func(c *gin.Context) {
		list := []*Broadcast{}
		h.broadcasts.Range(func(key, value interface{}) bool {
			list = append(list, value.(*Broadcast).Copy())
			return true
		})
		sort.SliceStable(list, func(i, j int) bool {
			return list[i].Participantduring10Min > list[j].Participantduring10Min
		})
		payload := Payload{
			Count: len(list),
			List:  list,
		}
		c.JSON(http.StatusOK, payload)
	}
	return f
}

func (h *Hub) PastBroadcast() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := strconv.Atoi(c.Query("uid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		collection := h.dClient.Database("livevup").Collection("broadcast")
		cursor, err := collection.Find(h.ctx, bson.D{{"uid", uid}},
			options.Find().SetSort(bson.D{{"livetime", -1}}),
			options.Find().SetLimit(10))
		defer cursor.Close(h.ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		broadcasts := []*ShortBroadcast{}
		for cursor.Next(h.ctx) {
			b := ShortBroadcast{}
			cursor.Decode(&b)
			broadcasts = append(broadcasts, &b)
		}
		c.JSON(http.StatusOK, broadcasts)
	}
}