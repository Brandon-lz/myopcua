package globaldata

var NodeIdWithValueToWrite = make(chan []NodeIdWithValueInput)
var NodeWriteResult = make(chan bool)


type NodeIdWithValueInput struct {
	NodeID   string  `json:"node_id" form:"node_id" binding:"required" example:"ns=2;i=2"`
	Value    interface{}  `json:"value" form:"value" binding:"required"`
	// DataType string `json:"data_type" form:"data_type" example:"Int"`
}

type NodeWriteResultOutput struct {
	NodeID   string  `json:"node_id"`
	Value    *string `json:"value"`
}