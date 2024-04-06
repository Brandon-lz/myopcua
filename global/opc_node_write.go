package globaldata

var NodeIdWithValueToWrite = make(chan []NodeIdWithValueInput)
var NodeWriteResult = make(chan bool)


type NodeIdWithValueInput struct {
	NodeName string  `json:"node_name" form:"node_name" binding:"required" example:"MyVariable"`  // 节点名称
	Value    interface{}  `json:"value" form:"value" binding:"required"`   // 要写入的值
	// DataType string `json:"data_type" form:"data_type" example:"Int"`
}

type NodeWriteResultOutput struct {
	NodeID   string  `json:"node_id"`
	Value    *string `json:"value"`
}