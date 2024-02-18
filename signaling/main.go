package main

import (
	"github.com/gin-gonic/gin"
	"signaling/handler"
)

func main() {
	r := gin.Default()
	r.GET("/", handler.Connection)
	r.Run(":8080")
}
