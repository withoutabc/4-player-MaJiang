package api

import (
	"github.com/gin-gonic/gin"
)

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
