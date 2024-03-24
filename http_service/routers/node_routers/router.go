package noderouters

import (
	globaldata "earth/global_data"
	"earth/http_service/core"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {
	nodegroup := router.Group("/opc-node")
	nodegroup.POST("/add-node-to-read", AddNodeToRead)
	nodegroup.GET("/get-node/:id", GetNode)
	nodegroup.DELETE("/delete-node/:id", DeleteNode)
}

// AddNodeToRead 路由 -------------------------------------------------------
// @Summary  AddNodeToRead 路由
// @Description  AddNodeToRead 路由
// @Tags     nodes
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param     data  body   AddNodeToReadRequest   true  "见下方JSON"
// @Success  200  {object}  AddNodeToReadResponse  "节点添加成功"
// @Router   /api/v1/opc-node/add-node-to-read [post]
func AddNodeToRead(c *gin.Context) {
	var req AddNodeToReadRequest
	core.BindParamAndValidate(c, &req)

	node := serviceAddNodeToRead(&req)

	out := core.SerializeDataAndValidate(OpcNodeOutput{node.NodeID, node.Name, node.DataType}, &OpcNodeOutput{},false)  // false不进行序列化，因为第一个参数类型就是目标类型

	// 响应
	core.SuccessHandler(c, AddNodeToReadResponse{
		Code:    200,
		Data:    out,
		Message: "节点添加成功",
	})

	// 为什么不是下面这样呢，这样子不够直观，而且也不方便将来做校验
	// SuccessHandler(c, "Create data_script successfully", data)
}

type AddNodeToReadRequest struct {
	Name     string  `json:"name" form:"name" binding:"required" example:"MyVariable"`
	NodeID   string  `json:"node-id" form:"node-id" binding:"required" example:"ns=2;i=2"`
	DataType *string `json:"data-type" form:"data-type" example:"Int32"`
}

type AddNodeToReadResponse struct {
	Code    int           `json:"code" example:"200"`
	Data    OpcNodeOutput `json:"data" `
	Message string        `json:"message" example:"节点添加成功"`
}

type OpcNodeOutput struct {
	Name     string `json:"name" example:"MyVariable"`
	NodeID   string `json:"node-id" example:"ns=2;s=MyVariable"`
	DataType string `json:"data-type" example:"Int32"`
}

func serviceAddNodeToRead(req *AddNodeToReadRequest) globaldata.OpcNode {
	// 配置新节点
	var node globaldata.OpcNode
	node.NodeID = req.NodeID
	node.Name = req.Name
	if req.DataType != nil {
		node.DataType = *req.DataType
	}
	if err := globaldata.SystemVars.AddNode(&node); err != nil {
		panic(core.NewKnownError(core.FailedToAddNode, req.NodeID, fmt.Sprintf("add node failed: %s", err.Error())))
	}
	globaldata.SystemVars.Save()
	return node
}

// GetNode 路由 -------------------------------------------------------
// @Summary  GetNode 路由
// @Description  GetNode 路由
// @Tags     nodes
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param     id  path  string  true  "节点ID"
// @Success  200  {object}  GetNodeResponse  "节点信息"
// @Router   /api/v1/opc-node/get-node/{id} [get]
func GetNode(c *gin.Context) {
	id := c.Param("id")
	// string to int64
	nodeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic(core.NewKnownError(http.StatusBadRequest, id, "Invalid node ID"))
	}

	node, value := serviceGetNode(nodeID)

	out := core.SerializeDataAndValidate(OpcNodeWithDataOutput{OpcNodeOutput{node.Name, node.NodeID, node.DataType}, value}, &OpcNodeWithDataOutput{},false)
	core.SuccessHandler(c, GetNodeResponse{
		Code:    200,
		Data:    out,
		Message: "节点信息",
	})
}

type GetNodeResponse struct {
	Code    int                   `json:"code" example:"200"`
	Data    OpcNodeWithDataOutput `json:"data" `
	Message string                `json:"message" example:"节点信息"`
}

type OpcNodeWithDataOutput struct {
	OpcNodeOutput
	Value string `json:"value" example:"123"`
}

func serviceGetNode(nodeID int64) (*globaldata.OpcNode, string) {
	node, err := globaldata.SystemVars.GetNode(nodeID)
	if err != nil {
		panic(core.NewKnownError(core.FailedToGetNode, nodeID, fmt.Sprintf("get node failed: %s", err.Error())))
	}

	value := ""
	if node.Value != nil {
		value = fmt.Sprintf("%v", node.Value)
	} else {
		value = "null"
	}
	return node, value
}

// DeleteNode 路由 -------------------------------------------------------
// @Summary  DeleteNode 路由
// @Description  DeleteNode 路由
// @Tags     nodes
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param     id  path  string  true  "节点ID"
// @Success  200  {object}  DeleteNodeResponse  "节点删除成功"
// @Router   /api/v1/opc-node/delete-node/{id} [delete]
func DeleteNode(c *gin.Context) {
	id := c.Param("id")
	// string to int64
	nodeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic(core.NewKnownError(http.StatusBadRequest, id, "Invalid node ID"))
	}
	ServiceDeleteNode(nodeID)
	core.SuccessHandler(c, DeleteNodeResponse{
		Code:    200,
		Data:    "",
		Message: "节点删除成功",
	})
}

type DeleteNodeResponse struct {
	Code    int    `json:"code" example:"200"`
	Data    string `json:"data" `
	Message string `json:"message" example:"节点删除成功"`
}

func ServiceDeleteNode(nodeID int64) {
	if err := globaldata.SystemVars.DeleteNode(nodeID); err != nil {
		panic(core.NewKnownError(core.FailedToDeleteNode, nodeID, "delete node failed"))
	}
	globaldata.SystemVars.Save()
}
