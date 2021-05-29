package worker

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Payload struct {
	Count int          `json:"count"`
	List  []*Broadcast `json:"list"`
}

func (h *Hub) List() gin.HandlerFunc {
	f := func(c *gin.Context) {
		list := []*Broadcast{}
		h.broadcasts.Range(func(key, value interface{}) bool {
			list = append(list, value.(*Broadcast).Copy())
			return true
		})
		payload := Payload{
			Count: len(list),
			List:  list,
		}
		c.JSON(http.StatusOK, payload)
	}
	return f
}