package api

import (
	"github.com/gin-gonic/gin"
)

// Rooms 定义房间列表
var Rooms = make(map[string]*Room)

func InitRouter() {

	router := gin.Default()

	u := router.Group("/user")
	{
		u.POST("/register", Register)
		u.POST("/login", Login)
	}

	router.GET("/ws", wsHandler)
	router.Run(":2022")
}
