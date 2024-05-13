package noderouters

import (
	"fmt"

	globaldata "github.com/Brandon-lz/myopcua/global"
	"github.com/Brandon-lz/myopcua/http_service/core"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {
	nodegroup := router.Group("/opc-node")
	nodegroup.POST("/add-node-to-read", AddNodeToRead)
	nodegroup.GET("/get-node/:id", GetNode)
	nodegroup.GET("/get-nodes", GetNodes)
	nodegroup.DELETE("/delete-node/:id", DeleteNode)
	nodegroup.PUT("/write-node-value", WriteNodeValue)
}

// AddNodeToRead 路由 -------------------------------------------------------
// @Summary  AddNodeToRead 路由 步骤**1**
// @Description  AddNodeToRead 路由
// @Tags     opc-nodes
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

	out := core.SerializeDataAndValidate(OpcNodeOutput{node.NodeID, node.Name, node.DataType}, &OpcNodeOutput{}, false) // false不进行序列化，因为第一个参数类型就是目标类型

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
	NodeID   string  `json:"node_id" form:"node_id" binding:"required" example:"ns=2;i=2"`
	DataType *string `json:"data_type" form:"data_type" example:"Int32"`
}

type AddNodeToReadResponse struct {
	Code    int           `json:"code" example:"200"`
	Data    OpcNodeOutput `json:"data" `
	Message string        `json:"message" example:"节点添加成功"`
}

type OpcNodeOutput struct {
	Name     string `json:"name" example:"MyVariable"`
	NodeID   string `json:"node_id" example:"ns=2;s=MyVariable"`
	DataType string `json:"data_type" example:"Int32"`
}

func serviceAddNodeToRead(req *AddNodeToReadRequest) globaldata.OpcNode {
	// 配置新节点
	var node globaldata.OpcNode
	node.NodeID = req.NodeID
	node.Name = req.Name
	if req.DataType != nil {
		node.DataType = *req.DataType
	}
	if err := globaldata.OPCNodeVars.AddNode(&node); err != nil {
		panic(core.NewKnownError(core.FailedToAddNode, req.NodeID, fmt.Sprintf("add node failed: %s", err.Error())))
	}
	globaldata.OPCNodeVars.Save()
	return node
}

// GetNode 路由 -------------------------------------------------------
// @Summary  GetNode 路由
// @Description  GetNode 路由
// @Tags     opc-nodes
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param     id  path  string  true  "节点ID"
// @Success  200  {object}  GetNodeResponse  "节点信息"
// @Router   /api/v1/opc-node/get-node/{id} [get]
func GetNode(c *gin.Context) {
	id := c.Param("id")
	// string to int64
	// nodeID, err := strconv.ParseInt(id, 10, 64)
	// if err != nil {
		// panic(core.NewKnownError(http.StatusBadRequest, id, "Invalid node ID"))
	// }

	node, value := serviceGetNode(id)

	out := core.SerializeDataAndValidate(OpcNodeWithDataOutput{OpcNodeOutput{node.Name, node.NodeID, node.DataType}, value}, &OpcNodeWithDataOutput{}, false)
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

func serviceGetNode(nodeID string) (*globaldata.OpcNode, string) {
	node, err := globaldata.OPCNodeVars.GetNodeByNodeId(nodeID)
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

// GetNodes 路由 ------------------------------------
// @Summary  GetNodes 路由
// @Description  GetNodes 路由
// @Tags     opc-nodes
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Success  200  {object}  GetNodesResponse  "节点列表"
// @Router   /api/v1/opc-node/get-nodes [get]
func GetNodes(c *gin.Context) {
	nodes := serviceGetNodes()
	out := core.SerializeDataAndValidate(nodes, &[]OpcNodeOutput{}, true)
	if out == nil {
		out = []OpcNodeOutput{}
	}
	core.SuccessHandler(c, GetNodesResponse{
		Code:    200,
		Data:    out,
		Message: "节点列表",
	})
}

type GetNodesResponse struct {
	Code    int           `json:"code" example:"200"`
	Data    []OpcNodeOutput `json:"data" `
	Message string        `json:"message" example:"节点列表"`
}

func serviceGetNodes() []OpcNodeOutput {
	var out []OpcNodeOutput

	for _,node := range globaldata.OPCNodeVars.CurrentNodes{
			out = append(out, OpcNodeOutput{node.Name, node.NodeID, node.DataType})
	}
	
	return out
}

// DeleteNode 路由 -------------------------------------------------------
// @Summary  DeleteNode 路由
// @Description  DeleteNode 路由
// @Tags     opc-nodes
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param     id  path  string  true  "节点ID"
// @Success  200  {object}  DeleteNodeResponse  "节点删除成功"
// @Router   /api/v1/opc-node/delete-node/{id} [delete]
func DeleteNode(c *gin.Context) {
	id := c.Param("id")

	// string to int64
	// nodeID, err := strconv.ParseInt(id, 10, 64)
	// if err != nil {
	// 	panic(core.NewKnownError(http.StatusBadRequest, id, "Invalid node ID"))
	// }
	ServiceDeleteNode(id)
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

func ServiceDeleteNode(nodeID string) {
	if err := globaldata.OPCNodeVars.DeleteNodeByNodeId(nodeID); err != nil {
		panic(core.NewKnownError(core.FailedToDeleteNode, nodeID, "delete node failed"))
	}
	globaldata.OPCNodeVars.Save()
}


// WriteNodeValue 路由 -------------------------------------------------------
// @Summary  WriteNodeValue 路由
// @Description  WriteNodeValue 路由
// @Description  参数定义
// @Description  ## 请求参数
// @Description  | 参数名称 | 类型 | 必填 | 描述 |
// @Description  | --- | --- | --- | --- |
// @Description  | data | list | 是 | 见下方JSON |
// @Description  ## data类型 定义
// @Description  | 字段 | 类型 | 是否必填 | 描述 |
// @Description  | --- | --- | --- | --- |
// @Description  | NodeName | string | 是 | 节点名称 |
// @Description  | Value | any | 是 | 写入值 |
// @Description  ## 参数示例
// @Description  ```json
// @Description  {
// @Description  "data": [
// @Description  {
// @Description  "NodeName": "MyVariable",
// @Description  "Value": 123
// @Description  },
// @Description  {
// @Description  "NodeName": "MyVariable2",
// @Description  "Value": "abc"
// @Description  }
// @Description  ]
// @Description  }
// @Description  ```
// @Description  ## 返回值定义
// @Description  | 字段 | 类型 | 描述 |
// @Description  | --- | --- | --- |
// @Description  | Code | int | 状态码 |
// @Description  | Message | string | 状态信息 |
// @Description  ## 返回值示例
// @Description  ```json
// @Description  {
// @Description  "Code": 200,
// @Description  "Message": "节点值写入完成"
// @Description  }
// @Description  ```
// @Tags     opc-nodes
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param     data  body   WriteNodeValueRequest   true  "见下方JSON"
// @Success  200  {object}  WriteNodeValueResponse  "节点值写入成功"
// @Router   /api/v1/opc-node/write-node-value [PUT]
func WriteNodeValue(c *gin.Context) {
	var req WriteNodeValueRequest
	core.BindParamAndValidate(c, &req)
	if len(req.Data) == 0 {
		panic(core.NewKnownError(core.FailedToWriteNodeValue, "", "data is empty"))
	}

	// check if node exist
	for _, data := range req.Data {
		_, err := globaldata.OPCNodeVars.GetNodeByName(data.NodeName)
		if err != nil {
			panic(core.NewKnownError(core.FailedToWriteNodeValue, data.NodeName, "node not found"))
		}
	}


	globaldata.OpcWriteLock.Lock()
	defer globaldata.OpcWriteLock.Unlock()

	globaldata.NodeIdWithValueToWrite <- req.Data

	writeRes := <-globaldata.NodeWriteResult
	if !writeRes {
		panic(core.NewKnownError(core.FailedToWriteNodeValue, "", "write node value failed"))
	}
	core.SuccessHandler(c, WriteNodeValueResponse{
		Code:    200,
		Message: "节点值写入完成",
	})

}

type WriteNodeValueRequest struct {
	Data []globaldata.NodeIdWithValueInput `json:"data" form:"data" binding:"required"`
}

type WriteNodeValueResponse struct {
	Code int `json:"code" example:"200"`
	// Data    []globaldata.NodeWriteResultOutput `json:"data"`
	Message string `json:"message" example:"节点值写入完成"`
}
