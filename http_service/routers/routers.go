package routers

import (
	globaldata "earth/global_data"
	"earth/http_service/core"
	alicerouters "earth/http_service/routers/alice_routers"
	smithrouters "earth/http_service/routers/smith_routers"
	"fmt"
	"net/http"
	"strconv"

	// "log"

	"github.com/gin-gonic/gin"
)



func RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", Root)
	router.GET("/ping", Ping)
	router.POST("/add-node-to-read", AddNodeToRead)
	router.GET("/get-node/:id", GetNode)
	router.DELETE("/delete-node/:id", DeleteNode)
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
	var node globaldata.OpcNode
	node.NodeID = req.NodeID
	node.Name = req.Name
	if req.DataType!= nil {
		node.DataType = *req.DataType
	}
	if err:=globaldata.SystemVars.AddNode(&node);err!=nil{
		panic(core.NewKnownError(core.FailedToAddNode,req.NodeID, fmt.Sprintf("add node failed: %s", err.Error())))
	}
	globaldata.SystemVars.Save()

	// 响应
	core.SuccessHandler(c, AddNodeToReadResponse{
		Code:    200,
		Data:    OpcNodeOutput{node.NodeID, node.Name, node.DataType},
		Message: "节点添加成功",
	})

	// 为什么不是下面这样呢，这样子不够直观，而且也不方便将来做校验
	// SuccessHandler(c, "Create data_script successfully", data)
}

type AddNodeToReadRequest struct {
	Name     string `json:"name" form:"name" validate:"required" example:"MyVariable"`
	NodeID string `json:"node-id" form:"node-id" validate:"required" example:"ns=2;i=2"`
	DataType *string `json:"data-type" form:"data-type" validate:"omitempty" example:"Int32"`
}

type AddNodeToReadResponse struct {
	Code    int    `json:"code" example:"200"`
	Data    OpcNodeOutput `json:"data" `
	Message string `json:"message" example:"节点添加成功"`
}

type OpcNodeOutput struct {
	Name     string `json:"name" example:"MyVariable"`
	NodeID   string `json:"node-id" example:"ns=2;s=MyVariable"`
	DataType string `json:"data-type" example:"Int32"`
}


// GetNode 路由
// @Summary  GetNode 路由
// @Description  GetNode 路由
// @Tags     default
// @Security BearerAuth	
// @Accept   json
// @Produce  json
// @Param     id  path  string  true  "节点ID"
// @Success  200  {object}  GetNodeResponse  "节点信息"
// @Router   /api/v1/get-node/{id} [get]
func GetNode(c *gin.Context) {
	id := c.Param("id")
	// string to int64
	nodeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic(core.NewKnownError(http.StatusBadRequest,id, "Invalid node ID"))
	}
	node, err := globaldata.SystemVars.GetNode(nodeID)
	if err != nil {
		panic(core.NewKnownError(core.FailedToGetNode,id, fmt.Sprintf("get node failed: %s", err.Error())))
	}

	value := ""
	if node.Value != nil {
		value = fmt.Sprintf("%v", node.Value)
	}else{
		value = "null"
	}

	core.SuccessHandler(c, GetNodeResponse{
		Code:    200,
		Data:    OpcNodeWithDataOutput{OpcNodeOutput{node.Name, node.NodeID, node.DataType}, value},
		Message: "节点信息",
	})
}


type GetNodeResponse struct {
	Code    int    `json:"code" example:"200"`
	Data    OpcNodeWithDataOutput `json:"data" `
	Message string `json:"message" example:"节点信息"`
}

type OpcNodeWithDataOutput struct {
	OpcNodeOutput
	Value string `json:"value" example:"123"`
}


// DeleteNode 路由
// @Summary  DeleteNode 路由
// @Description  DeleteNode 路由
// @Tags     default
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param     id  path  string  true  "节点ID"
// @Success  200  {object}  ApiResponse  "节点删除成功"
// @Router   /api/v1/delete-node/{id} [delete]
func DeleteNode(c *gin.Context) {
	id := c.Param("id")
	// string to int64
	nodeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic(core.NewKnownError(http.StatusBadRequest,id, "Invalid node ID"))
	}
	if err=globaldata.SystemVars.DeleteNode(nodeID);err!=nil{
		panic(core.NewKnownError(core.FailedToDeleteNode,id, "delete node failed"))
	}
	globaldata.SystemVars.Save()
	core.SuccessHandler(c, ApiResponse{
		Message: "节点删除成功",
	})
}

type DeleteNodeResponse struct {
	Code    int    `json:"code" example:"200"`
	Data    string `json:"data" `
	Message string `json:"message" example:"节点删除成功"`
}