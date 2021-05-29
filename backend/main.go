package main

import (
	"bililive/worker"
	"github.com/gin-gonic/gin"
)

func main() {
	hub := worker.Hub{}
	hub.Init()

	r := gin.Default()
	r.GET("/api/online", hub.List())
	r.Run()
}
