package webhookrouters

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Brandon-lz/myopcua/db/gen/model"
	"github.com/Brandon-lz/myopcua/db/gen/query"
	globaldata "github.com/Brandon-lz/myopcua/global_data"
	"github.com/Brandon-lz/myopcua/http_service/core"
	"github.com/Brandon-lz/myopcua/utils"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/webhook")
	group.POST("/example", WebHookExample)
	group.POST("/condition", CreateCondition)
	group.POST("/config", AddWebhookConfig)
	group.GET("/config/:id", GetWebhookConfigById)
	group.GET("/config-by-name/:name", GetWebhookConfigByName)
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
	var req WebHookExampleRequest
	core.BindParamAndValidate(c, &req)
	// fmt.Printf("req: %+v\n", req)
	utils.PrintDataAsJson(req)

	// 逻辑处理
	// 出参序列化以及校验
	core.SuccessHandler(c, WebHookExampleResponse{
		Code:    200,
		Data:    "webhook example",
		Message: "webhook example success",
	})
}

type WebHookExampleRequest struct {
	NodeName string `json:"node_name" form:"node_name" example:"MyVariable"`    // 节点名称
	Value    interface{} `json:"value" form:"value"`                   // 节点值 any类型
}

type WebHookExampleResponse struct {
	Code    int    `json:"code" example:"200"`
	Data    string `json:"data" example:"webhook example"`
	Message string `json:"message" example:"webhook example success"`
}

// AddWebhookConfig router -------------------------------------
// @Summary 配置一条新的webhook
// @Description # 配置一条新的webhook
// @Description
// @Description ## 例1：当节点node1值等于123时，发送通知到http://localhost:8080/api/v1/webhook/example
// @Description ```json
// @Description {
// @Description     "active": true,
// @Description     "name":"webhook1",
// @Description     "url": "http://localhost:8080/api/v1/webhook/example",
// @Description     "when": {
// @Description         "rule": {
// @Description 	 		"node_name": "node1",
// @Description              "type": "eq",
// @Description              "value": "123"
// @Description          }
// @Description      }
// @Description }
// @Description ```
// @Description 使用when字段会创建新的条件condition，并将其配置在这个webhook上
// @Description ## 例2：使用已经配置好的条件condition
// @Description ```json
// @Description {
// @Description    "active": true,
// @Description    "url": "http://localhost:8080/api/v1/webhook/example",
// @Description    "condition_id": 10
// @Description }
// @Description ```
// @Description ## 常见异常
// @Description - "code": 2007 代表数据重复，不能创建重复的webhook，具体重复了哪个字段，请看ConstraintName最后一个下划线后面的字段名
// @Description - "code": 400 "json: cannot unmarshal string into Go struct field AddWebhookConfigRequest.condition_id of type int64"  ： 看下body参数，数字类型传成了字符串
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
	if req.When != nil && req.ConditionId != nil {
		panic(core.NewKnownError(http.StatusBadRequest, nil, "when and condition_id cannot be both set"))
	}

	utils.PrintDataAsJson(req)

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
	Name        *string               `json:"name" form:"name" example:"webhook1"`                                            // webhook名称，可以为空
	Url         string                `json:"url" form:"url" binding:"required,url" example:"http://192.168.1.1:8800/notify"` // webhook地址
	Active      *bool                 `json:"active" form:"active" example:"true"`                                            // 是否激活，不传的话默认true
	When        *globaldata.Condition `json:"when" form:"when"`                                                               // 触发条件，为空时相当于通知所有数据变化
	ConditionId *int64                `json:"condition_id" form:"condition_id" example:"1"`                                   // 条件id，不传的话默认新增条件
	NeedNodeList []string `json:"need_node_list" form:"need_node_list" binding:"required" example:["node1,node2"]` // 需要的节点值列表，到时候会传参给webhook
}

type AddWebhookConfigResponse struct {
	Code    int               `json:"code" example:"200"`
	Data    WebHookConfigRead `json:"data" `
	Message string            `json:"message" example:"节点添加成功"`
}

type WebHookConfigRead struct {
	Id          int64     `json:"id" form:"id" validate:"required"`
	Name        string    `json:"name" form:"name" validate:"required"`
	Url         string    `json:"url" form:"url" validate:"required"`
	Active      bool      `json:"active" form:"active" validate:"required"`
	When        *string   `json:"when" form:"when" validate:"omitempty"`
	NeedNodeList []string `validate:"required"`
	ConditionId *int64    `json:"condition_id" form:"condition_id" validate:"omitempty"`
	CreatedAt   time.Time `json:"created_at" form:"created_at" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at" form:"updated_at" validate:"required"`
}

func ServiceAddWebhookConfig(req *AddWebhookConfigRequest) WebHookConfigRead {
	if req.Name == nil {
		req.Name = utils.Adr(uuid.New().String()[:6])
	}
	if req.Active == nil {
		req.Active = utils.Adr(true)
	}
	webhook, condition, needNodes := DalAddWebhookConfig(req)
	out := core.SerializeData(webhook, &WebHookConfigRead{}) // orm model -> out
	out.When = &condition.Condition
	if req.ConditionId != nil {
		out.ConditionId = req.ConditionId
	} else {
		out.ConditionId = &condition.ID
	}
	for _, needNode := range needNodes {
		out.NeedNodeList = append(out.NeedNodeList, needNode.NodeName)
	}
	return out
}

// {"name":"webhook1","url":"http://localhost:8080/api/v1/webhook/example","active":true,"when":{"and":null,"or":null,"rule":{"type":"eq","node_name":"node1","value":"123"}},"condition_id":null}
// {"name":"webhook1","url":"http://localhost:8080/api/v1/webhook/example","active":true,"when":{"and":null,"or":null,"rule":{"type":"eq","node_name":"node1","value":"123"}},"condition_id":null}
func DalAddWebhookConfig(req *AddWebhookConfigRequest) (*model.WebHook, *model.WebHookCondition,[]model.NeedNode) {
	var webhook model.WebHook
	var condition model.WebHookCondition
	var needNodes []model.NeedNode
	err := query.Q.Transaction(func(tx *query.Query) error {
		if req.When != nil {
			condition = model.WebHookCondition{Condition: utils.PrintDataAsJson(req.When)}
			err := tx.WebHookCondition.Create(&condition)
			if err != nil {
				slog.Error(utils.WrapError(err).Error())
				return err
			}
			// create webhook foreign
			webhook = core.SerializeData(req, &model.WebHook{}) // req -> orm model
			webhook.WebHookConditionRefer = &condition.ID
			err = tx.WebHook.Create(&webhook)
			if err != nil {
				slog.Error(utils.WrapError(err).Error())
				return err
			}
			// create need node foreign
			for _, nodeName := range req.NeedNodeList {
				var needNode = model.NeedNode{WebHookRefer: &webhook.ID, NodeName: nodeName}
				err = query.Q.NeedNode.Create(&needNode)
				if err != nil {
					slog.Error(utils.WrapError(err).Error())
					return err
				}
				needNodes = append(needNodes, needNode)
			}
		} else {
			webhook = core.SerializeData(req, &model.WebHook{}) // req -> orm model
			if req.ConditionId != nil {
				condition,err:=query.Q.WebHookCondition.Where(query.Q.WebHookCondition.ID.Eq(*req.ConditionId)).First()
				if err != nil || condition == nil {
					return core.NewKnownError(http.StatusBadRequest, err, "condition not exist")
				}
				webhook.WebHookConditionRefer = req.ConditionId
			}else{
				webhook.WebHookConditionRefer = nil
				webhook.Active = utils.Adr(false)
			}
			err := tx.WebHook.Create(&webhook)
			if err != nil {
				slog.Error(utils.WrapError(err).Error())
				return err
			}
			// create need node foreign
			for _, nodeName := range req.NeedNodeList {
				var needNode = model.NeedNode{WebHookRefer: &webhook.ID, NodeName: nodeName}
				err = query.Q.NeedNode.Create(&needNode)
				if err != nil {
					slog.Error(utils.WrapError(err).Error())
					return err
				}
				needNodes = append(needNodes, needNode)
			}
		}
		return nil
	})
	if err != nil {
		slog.Error(utils.WrapError(err).Error())
		sqlErr := err.(*pgconn.PgError)
		panic(core.NewKnownError(core.FieldNotUnique, err, sqlErr.Message))
	}

	fmt.Printf("webhook: %+v\n", webhook)
	fmt.Printf("condition: %+v\n", condition)

	return &webhook, &condition, needNodes
}

// GetWebhookConfigById router -------------------------------
// @Summary 根据id获取webhook配置
// @Description 根据id获取webhook配置
// @Tags Webhook
// @Accept  json
// @Produce  json
// @Param id path string true "webhook id"
// @Success 200 {object} GetWebhookConfigByIdResponse
// @Router /api/v1/webhook/config/{id} [get]
func GetWebhookConfigById(c *gin.Context) {
	// 入参校验
	id := c.Param("id")
	if id == "" {
		panic(core.NewKnownError(http.StatusBadRequest, nil, "id is empty"))
	}

	out := ServiceGetWebhookConfigById(id)
	slog.Debug(fmt.Sprintf("out: %+v", out.NeedNodeList))
	core.ValidateSchema(out)

	core.SuccessHandler(c, GetWebhookConfigByIdResponse{
		Code:    200,
		Data:    out,
		Message: "Webhook configuration get successfully",
	})
}

type GetWebhookConfigByIdResponse struct {
	Code    int               `json:"code" example:"200"`
	Data    WebHookConfigRead `json:"data" `
	Message string            `json:"message" example:"Webhook configuration get successfully"`
}

func ServiceGetWebhookConfigById(id string) WebHookConfigRead {
	// 逻辑处理
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic(core.NewKnownError(http.StatusBadRequest, err, "id is not int"))
	}
	return GetWebhookConfigByIdFromDB(intId)
}

func GetWebhookConfigByIdFromDB(id int64) WebHookConfigRead {
	webhook, condition := DalGetWebhookConfigById(id)
	// 出参序列化以及校验
	out := core.SerializeData(webhook, &WebHookConfigRead{})
	if condition != nil {
		out.When = &condition.Condition
	}
	return out
}

func GetAllWebhookConfigFromDB() []WebHookConfigRead {
	var out []WebHookConfigRead
	tuples, err := globaldata.DalGetAllWebhookConfig()
	if err != nil {
		slog.Error(utils.WrapError(err).Error())
		panic(core.NewKnownError(core.EntityNotFound, err, "webhook not found"))
	}
	for _, tuple := range tuples {
		var webhook WebHookConfigRead = core.SerializeData(tuple.Webhook, &WebHookConfigRead{})
		if tuple.Condition != nil {
			webhook.Id = tuple.Condition.ID
			webhook.When = &tuple.Condition.Condition
		}
		out = append(out, webhook)
	}
	return out
}

func DalGetWebhookConfigById(id int64) (*model.WebHook, *model.WebHookCondition) {
	var webhook *model.WebHook
	var err error
	q := query.Q.WebHook
	webhook, err = q.Where(q.ID.Eq(id)).First()
	if err != nil {
		slog.Error(utils.WrapError(err).Error())
		panic(core.NewKnownError(core.EntityNotFound, err, "webhook not found"))
	}
	if webhook.WebHookConditionRefer == nil {
		return webhook, nil
	}
	// gen foreign key condition not support yet, deal with handly
	u := query.Q.WebHookCondition
	condition, err := u.Where(u.ID.Eq(*webhook.WebHookConditionRefer)).First()
	if err != nil {
		slog.Error(utils.WrapError(err).Error())
		panic(core.NewKnownError(core.EntityNotFound, err, "condition not found"))
	}
	return webhook, condition
}

// GetWebhookConfigByName router ------------------------------
// @Summary 根据名称获取webhook配置
// @Description 根据名称获取webhook配置
// @Tags Webhook
// @Accept  json
// @Produce  json
// @Param name path string true "webhook名称"
// @Success 200 {object} GetWebhookConfigByNameResponse
// @Router /api/v1/webhook/config-by-name/{name} [get]
func GetWebhookConfigByName(c *gin.Context) {
	// 入参校验
	name := c.Param("name")
	if name == "" {
		panic(core.NewKnownError(http.StatusBadRequest, nil, "name is empty"))
	}

	// 逻辑处理
	webhook := DalGetWebhookConfigByName(name)

	// 出参序列化以及校验
	out := core.SerializeData(webhook, &WebHookConfigRead{})
	core.ValidateSchema(out)

	core.SuccessHandler(c, GetWebhookConfigByNameResponse{
		Code:    200,
		Data:    out,
		Message: "Webhook configuration get successfully",
	})
}

type GetWebhookConfigByNameResponse struct {
	Code    int               `json:"code" example:"200"`
	Data    WebHookConfigRead `json:"data" `
	Message string            `json:"message" example:"Webhook configuration get successfully"`
}

func DalGetWebhookConfigByName(name string) *model.WebHook {
	var webhook *model.WebHook
	var err error
	q := query.Q.WebHook
	webhook, err = q.Where(q.Name.Eq(name)).First()
	if err != nil {
		panic(core.NewKnownError(core.EntityNotFound, err, "webhook not found"))
	}

	fmt.Printf("webhook: %+v\n", webhook)

	return webhook
}

// CreateCondition router -------------------------------------
// @Summary 创建触发条件
// @Description # 参数说明
// @Description ## 请求参数
// @Description | 参数名称 | 类型 | 必填 | 描述 |
// @Description | --- | --- | --- | --- |
// @Description | and | list中嵌套本参数 | 否 | 规则列表，逻辑与 |
// @Description | or | list中嵌套本参数 | 否 | 规则列表，逻辑或 |
// @Description | rule | Rule | 否 | 规则 |
// @Description ## Rule类型 定义
// @Description | 字段 | 类型 | 是否必填 | 描述 |
// @Description | --- | --- | --- | --- |
// @Description | node_name | string | 是 | 节点名称 |
// @Description | type | string | 是 | 规则类型，支持eq ne gt lt all-time in not-in |
// @Description | value | any | 是 | 比对值 |
// @Description ## 参数示例1 : 当节点MyVariable大于123时触发
// @Description ```json
// @Description {
// @Description     "rule": {
// @Description         "node_name": "MyVariable",
// @Description         "type": "gt",
// @Description         "value": 123
// @Description     }
// @Description }
// @Description ```
// @Description ## 参数示例2 : 当节点node1等于在["abc","def"]，并且节点node2等于123时触发
// @Description ```json
// @Description {
// @Description     "and": [
// @Description         {
// @Description             "rule": {
// @Description                 "node_name": "node1",
// @Description                 "type": "in",
// @Description                 "value": [
// @Description                     "abc",
// @Description                     "def",
// @Description             }
// @Description         },
// @Description         {
// @Description             "rule": {
// @Description                 "node_name": "node2",
// @Description                 "type": "eq",
// @Description                 "value": 123
// @Description             }
// @Description         }
// @Description     ]
// @Description }
// @Description ```
// @Description *注意：Condition是嵌套类型，Condition包含and，or，rule，所以and里面可以嵌套and。。。无限嵌套*
// @Tags Webhook
// @Accept  json
// @Produce  json
// @Param body body CreateConditionRequest true "创建条件"
// @Success 200 {object} CreateConditionResponse
// @Router /api/v1/webhook/condition [post]
func CreateCondition(c *gin.Context) {
	// 入参校验
	var req CreateConditionRequest
	core.BindParamAndValidate(c, &req)

	// log.Logger.Debug("req: ", utils.PrintDataAsJson(req.Rule.Value))

	// 逻辑处理
	condition := ServiceCreateCondition(&req)

	// 出参序列化以及校验
	out := core.SerializeData(condition, &WebHookConditionRead{})
	core.ValidateSchema(out)

	core.SuccessHandler(c, CreateConditionResponse{
		Code:    200,
		Data:    out,
		Message: "Condition created successfully",
	})
}

type CreateConditionRequest struct {
	And  []globaldata.Condition `json:"and" form:"and"`   // 规则列表，逻辑与
	Or   []globaldata.Condition `json:"or" form:"or"`     // 规则列表，逻辑或
	Rule *globaldata.Rule       `json:"rule" form:"rule"` // 规则
}

type CreateConditionResponse struct {
	Code    int                  `json:"code" example:"200"`
	Data    WebHookConditionRead `json:"data" `
	Message string               `json:"message" example:"Condition created successfully"`
}

type WebHookConditionRead struct {
	Id        int64     `json:"id" form:"id" validate:"required"`
	Condition string    `json:"condition" form:"condition" validate:"required"`
	CreatedAt time.Time `json:"created_at" form:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at" validate:"required"`
}

func ServiceCreateCondition(req *CreateConditionRequest) WebHookConditionRead {
	// check params
	if req.And == nil && req.Or == nil && req.Rule == nil {
		panic(core.NewKnownError(http.StatusBadRequest, nil, "params is empty"))
	}
	// 解析即校验???
	
	// to db
	condition := DalCreateCondition(req)
	out := core.SerializeData(condition, &WebHookConditionRead{}) // orm model -> out
	return out
}

func DalCreateCondition(req *CreateConditionRequest) *model.WebHookCondition {
	var condition model.WebHookCondition
	conditionData := map[string]string{
		"condition": utils.PrintDataAsJson(req),
	}
	err := query.Q.Transaction(func(tx *query.Query) error {
		condition = core.SerializeData(conditionData, &model.WebHookCondition{}) // req -> orm model
		err := tx.WebHookCondition.Create(&condition)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		sqlErr := err.(*pgconn.PgError)
		panic(core.NewKnownError(core.FieldNotUnique, err, sqlErr.Message))
	}

	return &condition
}
