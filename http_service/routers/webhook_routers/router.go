package webhookrouters

import (
	"earth/http_service/core"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/webhook")
	group.POST("/config", AddWebhookConfig)
	// group.GET("/config/:id", GetWebhookConfig)
	// group.GET("/configs", GetAllWebhookConfigs)
	// group.DELETE("/config/:id", DeleteWebhookConfig)
}

// AddWebhookConfig router -------------------------------------
// @Summary Add a new webhook configuration
// @Description Add a new webhook configuration
// @Tags Webhook
// @Accept  json
// @Produce  json
// @Param body body AddWebhookConfigRequest true "Webhook configuration"
// @Success 200 {object} AddWebhookConfigResponse
// @Router /api/v1/webhook/config [post]
func AddWebhookConfig(c *gin.Context) {
	// 入参校验
	var req AddWebhookConfigRequest
	core.BindParamAndValidate(c, &req)
	fmt.Printf("req: %+v\n", req)
	
	// 逻辑处理
	webhook,err:=ServiceAddWebhookConfig(&req)
	if err!=nil{
		panic(err)
	}

	// 出参序列化以及校验
	out:=core.SerializeDataAndValidate(*webhook, &WebHookConfigRead{},false)   // false代表只校验字段但是不做序列化，因为这里的webhook变量已经是目标类型了

	core.SuccessHandler(c, AddWebhookConfigResponse{
		Code:    200,
		Data:    out,
		Message: "Webhook configuration added successfully",
	})
}

type AddWebhookConfigRequest struct {
	Id     int     `json:"id" form:"id" binding:"required,gt=10" example:"1"`
	Name   *string `json:"name" form:"name" example:"webhook1"` // 可以为空 要用*string
	Url    string  `json:"url" form:"url" binding:"required"`
	Active *bool   `json:"active" form:"active"`
}

type AddWebhookConfigResponse struct {
	Code    int               `json:"code" example:"200"`
	Data    WebHookConfigRead `json:"data" `
	Message string            `json:"message" example:"节点添加成功"`
}

type WebHookConfigRead struct {
	Id     int    `json:"id" form:"id" validate:"required"`
	Name   string `json:"name" form:"name" validate:"required"`
	Url    string `json:"url" form:"url" validate:"required"`
	Active bool   `json:"active" form:"active" validate:"required"`
}


func ServiceAddWebhookConfig(req *AddWebhookConfigRequest) (*WebHookConfigRead, error) {
	var resp WebHookConfigRead
	resp.Id = req.Id
	if req.Name == nil {
		resp.Name = "webhook"
	} else {
		resp.Name = uuid.New().String()[:6]
	}
	resp.Url = req.Url
	if req.Active == nil {
		resp.Active = true
	} else {
		resp.Active = *req.Active
	}

	return &resp, nil
}

