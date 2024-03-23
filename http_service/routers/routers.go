package routers

import (
	"earth/http_service/core"
	alicerouters "earth/http_service/routers/alice_routers"
	smithrouters "earth/http_service/routers/smith_routers"

	// "earth/http_service/utils"

	"github.com/gin-gonic/gin"
)



func RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", Root)
	router.GET("/ping", Ping)
	router.POST("/add-node-to-read", AddNodeToRead)
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



// AddNodeToRead 路由
// @Summary  AddNodeToRead 路由
// @Description  AddNodeToRead 路由
// @Tags     default
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param     data  body   AddNodeToReadRequest   true  "见下方JSON"
// @Success  200  {object}  AddNodeToReadResponse  "节点添加成功"
// @Router   /api/v1/add-node-to-read [post]
func AddNodeToRead(c *gin.Context) {
	var req AddNodeToReadRequest
	core.BindParamAndValidate(c, &req)
	// 配置新节点
	// globaldata.SystemVars.AddNode()

	// 响应
	core.SuccessHandler(c, AddNodeToReadResponse{
		Code:    200,
		Data:    "",
		Message: "节点添加成功",
	})

	// 为什么不是下面这样呢，这样子不够直观，而且也不方便将来做校验
	// SuccessHandler(c, "Create data_script successfully", data)
}

// type OpcNode struct {
// 	NodeID   string
// 	Name     string
// 	DataType string
// 	Value    interface{}
// }

type AddNodeToReadRequest struct {
	Name     string `json:"name" form:"name" validate:"required" example:"MyVariable"`
	NodeID string `json:"node-id" form:"node-id" validate:"required" example:"ns=2;s=MyVariable"`
	DataType string `json:"data-type" form:"data-type" validate:"omitempty" example:"Int32"`
}

type AddNodeToReadResponse struct {
	Code    int    `json:"code" example:"200"`
	Data    string `json:"data" example:""`
	Message string `json:"message" example:"节点添加成功"`
}