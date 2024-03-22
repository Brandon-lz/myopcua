package smithrouters

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {
	smithgroup := router.Group("/smith")
	smithgroup.GET("/ping", MyPing)
}

func MyPing (c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

