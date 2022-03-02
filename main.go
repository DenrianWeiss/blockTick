package main

import (
	"github.com/DenrianWeiss/blockTick/config"
	"github.com/DenrianWeiss/blockTick/handler"
	"github.com/DenrianWeiss/blockTick/service"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	service.InitAll()
	r := gin.Default()
	r.GET("/", handler.MainPageHandler)
	r.POST("/api/v1/", handler.ReceiveTransaction)
}
