package routers

import (
	alicerouters "earth/http_service/routers/alice_routers"
	noderouters "earth/http_service/routers/node_routers"
	smithrouters "earth/http_service/routers/smith_routers"
	webhookrouters "earth/http_service/routers/webhook_routers"

	// "log"

	"github.com/gin-gonic/gin"
)



func RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", Root)
	router.GET("/ping", Ping)
	noderouters.RegisterRoutes(router)
	webhookrouters.RegisterRoutes(router)
	smithrouters.RegisterRoutes(router)
	alicerouters.RegisterRoutes(router)
}


// 根路由
// @Summary  根路由
// @Description  根路由
// @Tags     default
// @Accept   json
// @Produce  json
// @Success  200  {object}  ApiResponse  "欢迎使用OPC-UA OpenAPI"
// @Router   /api/v1/ [get]
func Root(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "欢迎使用OPC-UA OpenAPI",
	})
}

type ApiResponse struct {
	Message string `json:"message" example:"欢迎使用OPC-UA OpenAPI"`
}

// Ping 路由
// @Summary  Ping 路由
// @Description  Ping 路由
// @Tags     default
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Success  200  {object}  ApiResponse  "pong"
// @Router   /api/v1/ping [get]
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

