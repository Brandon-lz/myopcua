package alicerouters


import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {
	alicegroup := router.Group("/alice")
	alicegroup.GET("/ping", MyPing)
}


type ApiResponse struct {
	Message string `json:"message" example:"pong"`
}

// Ping 路由
// @Summary  Ping 路由
// @Description  Ping 路由
// @Tags     alice的接口
// @Accept   json
// @Produce  json
// @Success  200  {object}  ApiResponse  "pong"
// @Router   /api/v1/alice/ping [get]
func MyPing (c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

