package worker

import (
	"fmt"
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
	Uname            string    `json:"uname" bson:"uname"`
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

type RankItem struct {
	UID            int64   `json:"uid" bson:"uid"`
	Uname          string  `json:"uname" bson:"uname"`
	Duration       float64 `json:"duration" bson:"duration"`
	Income         float64 `json:"income" bson:"income"`
	DanmuCount     int64   `json:"danmuCount" bson:"danmuCount"`
	AvgPaidUser    float64 `json:"avgPaidUser" bson:"avgPaidUser"`
	AvgParticipant float64 `json:"avgParticipant" bson:"avgParticipant"`
	AvgViewership  float64 `json:"avgViewership" bson:"avgViewership"`
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

func (h *Hub) Rank() gin.HandlerFunc {
	return func(c *gin.Context) {
		sortBy := c.Query("sortBy")
		collection := h.dClient.Database("livevup").Collection("weekly")

		ops := []*options.FindOptions{options.Find().SetLimit(50)}
		switch sortBy {
		case "duration":
			ops = append(ops, options.Find().SetSort(bson.D{{"duration", -1}}))
		case "viewership":
			ops = append(ops, options.Find().SetSort(bson.D{{"avgViewership", -1}}))
		case "paid":
			ops = append(ops, options.Find().SetSort(bson.D{{"avgPaidUser", -1}}))
		default: // "income" or anything else
			ops = append(ops, options.Find().SetSort(bson.D{{"income", -1}}))
		}

		cursor, err := collection.Find(h.ctx, bson.D{}, ops...)
		defer cursor.Close(h.ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			fmt.Printf("worker/api: %v\n", err)
			return
		}

		res := []*RankItem{}
		for cursor.Next(h.ctx) {
			r := RankItem{}
			cursor.Decode(&r)
			res = append(res , &r)
		}
		c.JSON(http.StatusOK, res)
	}
}