package worker

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

type Payload struct {
	Count int          `json:"count"`
	List  []*Broadcast `json:"list"`
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