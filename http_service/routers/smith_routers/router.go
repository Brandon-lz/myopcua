package smithrouters

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {
	smithgroup := router.Group("/smith")
	smithgroup.GET("/ping", MyPing)
}

type ApiResponse struct {
	Message string `json:"message" example:"pong"`
}



func MyPing (c *gin.Context) {
	// Ping 路由
	// @Summary  Ping 路由
	// @Description  Ping 路由
	// @Tags     smith的接口
	// @Accept   json
	// @Produce  json
	// @Success  200  {object}  ApiResponse  "pong"
	// @Router   /api/v1/smith/ping [get]
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

