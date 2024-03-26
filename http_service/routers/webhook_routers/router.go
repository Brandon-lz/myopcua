package webhookrouters

import (
	"fmt"
	"time"

	"github.com/Brandon-lz/myopcua/db/gen/model"
	"github.com/Brandon-lz/myopcua/db/gen/query"
	"github.com/Brandon-lz/myopcua/http_service/core"
	"github.com/Brandon-lz/myopcua/utils"
	"github.com/jackc/pgx/v5/pgconn"

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
	webhook := ServiceAddWebhookConfig(&req)

	// 出参序列化以及校验
	out := core.SerializeData(webhook, &WebHookConfigRead{})
	core.ValidateSchema(out)

	core.SuccessHandler(c, AddWebhookConfigResponse{
		Code:    200,
		Data:    out,
		Message: "Webhook configuration added successfully",
	})
}

// GetWebhookConfig router 参数定义，字段描述放在字段后面
type AddWebhookConfigRequest struct {
	Name   *string `json:"name" form:"name" example:"webhook1"` // webhook名称，可以为空
	Url    string  `json:"url" form:"url" binding:"url" example:"http://192.168.1.1:8800/notify"`   // webhook地址
	Active *bool   `json:"active" form:"active" example:"true"`                // 是否激活，不传的话默认true
}

type AddWebhookConfigResponse struct {
	Code    int               `json:"code" example:"200"`
	Data    WebHookConfigRead `json:"data" `
	Message string            `json:"message" example:"节点添加成功"`
}

// type WebHook struct {
// 	ID        int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
// 	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
// 	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
// 	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
// 	Name      string         `gorm:"column:name;not null;comment:webhook名称" json:"name"`    // webhook名称
// 	URL       string         `gorm:"column:url;not null;comment:url地址" json:"url"`          // url地址
// 	Active    bool           `gorm:"column:active;default:true;comment:是否激活" json:"active"` // 是否激活
// }

type WebHookConfigRead struct {
	Id        int       `json:"id" form:"id" validate:"required"`
	Name      string    `json:"name" form:"name" validate:"required"`
	Url       string    `json:"url" form:"url" validate:"required"`
	Active    bool      `json:"active" form:"active" validate:"required"`
	CreatedAt time.Time `json:"created_at" form:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at" validate:"required"`
}

func ServiceAddWebhookConfig(req *AddWebhookConfigRequest) WebHookConfigRead {
	if req.Name == nil {
		req.Name = utils.Adr(uuid.New().String()[:6])
	}
	if req.Active == nil {
		req.Active = utils.Adr(true)
	}
	webhook := DalAddWebhookConfig(req)
	out := core.SerializeData(webhook, &WebHookConfigRead{}) // orm model -> out
	return out
}

func DalAddWebhookConfig(req *AddWebhookConfigRequest) *model.WebHook {

	webhook := core.SerializeData(req, &model.WebHook{}) // req -> orm model

	fmt.Printf("webhook: %+v\n", webhook)

	err := query.Q.WebHook.Create(&webhook) // auto print error, no need to log by hand
	if err != nil {
		sqlErr := err.(*pgconn.PgError)
		panic(core.NewKnownError(core.FieldNotUnique, err, sqlErr.Message))
	}
	return &webhook
}
