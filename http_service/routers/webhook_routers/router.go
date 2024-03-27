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
	group.POST("/example", WebHookExample)
	// group.PUT("/config/:id", UpdateWebhookConfig)
	// group.GET("/config/:id", GetWebhookConfig)
	// group.GET("/configs", GetAllWebhookConfigs)
	// group.DELETE("/config/:id", DeleteWebhookConfig)
}


// WebHookExample router ---------------------------------------
// @Summary webhook示例
// @Description webhook示例
// @Tags Webhook
// @Accept  json
// @Produce  json
// @Param body body WebHookExampleRequest true "Webhook example"
// @Success 200 {object} WebHookExampleResponse
// @Router /api/v1/webhook/example [POST]
func WebHookExample(c *gin.Context) {
	// 入参校验
	var req AddWebhookConfigRequest
	core.BindParamAndValidate(c, &req)
	fmt.Printf("req: %+v\n", req)

	// 逻辑处理
	// 出参序列化以及校验
	core.SuccessHandler(c, WebHookExampleResponse{
		Code:    200,
		Data:    "webhook example",
		Message: "webhook example success",
	})
}

type WebHookExampleRequest struct {
	NodeName string `json:"node_name" form:"node_name" example:"MyVariable"` // 节点名称
	NodeId string `json:"node_id" form:"node_id" example:"ns=1;s=MyVariable"` // 节点id
	Value string `json:"value" form:"value" example:"123"` // 入参示例
}

type WebHookExampleResponse struct {
	Code    int    `json:"code" example:"200"`
	Data    string `json:"data" example:"webhook example"`
	Message string `json:"message" example:"webhook example success"`
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
	When  *Condition   `json:"when" form:"when"` // 触发条件，为空时相当于通知所有数据变化
}

// type When struct {
// 	And []Condition `json:"and" form:"and" example:"[{\"and\":[{\"type\":\"eq\",\"key\":\"tag1\",\"value\":\"value1\"}]}]"` // 规则列表，逻辑与
// 	Or  []Condition  `json:"or" form:"or" example:"[{\"or\":[{\"type\":\"eq\",\"key\":\"tag1\",\"value\":\"value1\"}]}]"`  // 规则列表，逻辑或
// 	Rule *Rule `json:"rule" form:"rule" example:"{\"type\":\"eq\",\"key\":\"tag1\",\"value\":\"value1\"}"` // 规则
// }

// type And struct {
// 	And []And `json:"and" form:"and" example:"[{\"type\":\"eq\",\"key\":\"tag1\",\"value\":\"value1\"}]"` // 规则列表
// 	Or  []Or `json:"or" form:"or" example:"[{\"type\":\"eq\",\"key\":\"tag1\",\"value\":\"value1\"}]"`  // 规则列表
// 	Rule *Rule `json:"rule" form:"rule" example:"{\"type\":\"eq\",\"key\":\"tag1\",\"value\":\"value1\"}"` // 规则
// }

type Condition struct {
	And []Condition `json:"and" form:"and"` // 规则列表，逻辑与
	Or  []Condition `json:"or" form:"or"`  // 规则列表，逻辑或
	Rule *Rule `json:"rule" form:"rule"` // 规则
}

type Rule struct {
	Type  string `json:"type" form:"type" example:"eq"` // 规则类型
	NodeName string `json:"node_name" form:"node_name" example:"MyVariable"` // 节点名称
	Value string `json:"value" form:"value" example:"123"` // 规则value
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
	var webhook model.WebHook
	var condition model.WebHookCondition
	err:= query.Q.Transaction(func(tx *query.Query) error {
		if req.When != nil {
			condition = model.WebHookCondition{Condition: utils.PrintMapAsJson(req.When)}
			err := tx.WebHookCondition.Create(&condition)
			if err != nil {
				return err
			}
			webhook = core.SerializeData(req, &model.WebHook{}) // req -> orm model
			webhook.WebHookConditionRefer = condition.ID
			err = tx.WebHook.Create(&webhook)
			if err != nil {
				return err
			}
		}else{
			webhook = core.SerializeData(req, &model.WebHook{}) // req -> orm model
			err := tx.WebHook.Create(&webhook)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		sqlErr := err.(*pgconn.PgError)
		panic(core.NewKnownError(core.FieldNotUnique, err, sqlErr.Message))
	}

	fmt.Printf("webhook: %+v\n", webhook)
	fmt.Printf("condition: %+v\n", condition)

	return &webhook
}


// func Condition