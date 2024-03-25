package webhookrouters

import (
	"context"
	"earth/db/gen/model"
	"earth/db/gen/query"
	"earth/http_service/core"
	"earth/utils"
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
// @Summary 配置一条新的webhook
// @Description 配置一条新的webhook
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
	// webhook,err:=ServiceAddWebhookConfig(&req)
	// if err!=nil{
	// 	panic(err)
	// }

	if req.Name == nil {
		// name := 
		req.Name = utils.Adr(uuid.New().String()[:6])
	}
	if req.Active == nil{
		req.Active = utils.Adr(true)
	}

	webhook := core.SerializeData(&req, &model.WebHook{})
	fmt.Println(11111111,webhook)

	q := query.WebHook
	ctx := context.Background()
	err:= q.WithContext(ctx).Create(&webhook)
	if err!=nil{
		panic(err)
	}
	fmt.Printf("%+v",webhook)

	// 出参序列化以及校验
	out:=core.SerializeData(webhook,&WebHookConfigRead{})
	// out:=core.SerializeDataAndValidate(webhook, &WebHookConfigRead{},true)   // false代表只校验字段但是不做序列化，因为这里的webhook变量已经是目标类型了
	core.ValidateSchema(out)

	fmt.Printf("%+v",out)

	core.SuccessHandler(c, AddWebhookConfigResponse{
		Code:    200,
		Data:    out,
		Message: "Webhook configuration added successfully",
	})
}

type AddWebhookConfigRequest struct {
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
