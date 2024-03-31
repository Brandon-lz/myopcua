package webhookrouters

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Brandon-lz/myopcua/db/gen/model"
	"github.com/Brandon-lz/myopcua/db/gen/query"
	globaldata "github.com/Brandon-lz/myopcua/global"
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
	group.PUT("/condition/:id", UpdateCondition)
	group.DELETE("/condition/:id", DeleteCondition)
	group.GET("/condition/:id", GetConditionById)
	group.GET("/conditions", GetAllConditionsByPage)
	group.POST("", AddWebhookConfig)
	group.GET("/:id", GetWebhookConfigById)
	group.GET("", GetWebhookConfigByName)
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
	slog.Debug(fmt.Sprintf("received webhook example request: %+v\n", req))

	// 逻辑处理
	// 出参序列化以及校验
	core.SuccessHandler(c, WebHookExampleResponse{
		Code:    200,
		Data:    "webhook example",
		Message: "webhook example success",
	})
}

type WebHookExampleRequest struct {
	Values map[string]interface{} `json:"values" form:"values" binding:"required"` // 节点值 any类型
}

type WebHookExampleResponse struct {
	Code    int    `json:"code" example:"200"`
	Data    string `json:"data" example:"webhook example"`
	Message string `json:"message" example:"webhook example success"`
}

// AddWebhookConfig router -------------------------------------
// @Summary 配置一条新的webhook
// @Description # 配置一条新的webhook
// @Description ## 说明
// @Description 该接口用于配置一条新的webhook，并通过when字段或condition_id字段配置触发条件，当条件满足时会触发webhook，并将所需要的数据传给webhook url接口。
// @Description ## 请求参数
// @Description
// @Description - name：webhook名称，为空时系统会自动生成一个uuid
// @Description - url：webhook地址（POST请求），必填，当条件满足时会调用该url，并挂载数据到body中
// @Description - active：是否激活，不传的话默认true
// @Description - when：条件，里面是一个json，具体的格式看/api/v1/webhook/condition接口说明，when其实就是一个condition条件。该字段和condition_id字段必须传一个。
// @Description - condition_id：条件id，将条件condition配置到webhook上，该字段和when字段必须传一个。
// @Description - need_node_list：需要的节点值列表，条件触发时会传参给webhook
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
// @Router /api/v1/webhook [post]
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
	Name         *string               `json:"name" form:"name" example:"webhook1"`                                            // webhook名称，可以为空
	Url          string                `json:"url" form:"url" binding:"required,url" example:"http://192.168.1.1:8800/notify"` // webhook地址
	Active       *bool                 `json:"active" form:"active" example:"true"`                                            // 是否激活，不传的话默认true
	When         *globaldata.Condition `json:"when" form:"when"`                                                               // 触发条件，为空时相当于通知所有数据变化
	ConditionId  *int64                `json:"condition_id" form:"condition_id" example:"1"`                                   // 条件id，不传的话默认新增条件
	NeedNodeList []string              `json:"need_node_list" form:"need_node_list" binding:"required"`                        // 需要的节点值列表，到时候会传参给webhook
} // Todo: NeedNodeList 解除required限制，因为有时候不需要汇报数据

type AddWebhookConfigResponse struct {
	Code    int               `json:"code" example:"200"`
	Data    WebHookConfigRead `json:"data"`
	Message string            `json:"message" example:"节点添加成功"`
}

type WebHookConfigRead struct {
	Id           int64     `json:"id" form:"id" validate:"required"`
	Name         string    `json:"name" form:"name" validate:"required"`
	Url          string    `json:"url" form:"url" validate:"required"`
	Active       bool      `json:"active" form:"active" validate:"required"`
	When         *string   `json:"when" form:"when" validate:"omitempty"`
	NeedNodeList []string  `validate:"required"`
	ConditionId  *int64    `json:"condition_id" form:"condition_id" validate:"omitempty"`
	CreatedAt    time.Time `json:"created_at" form:"created_at" validate:"required"`
	UpdatedAt    time.Time `json:"updated_at" form:"updated_at" validate:"required"`
}

func ServiceAddWebhookConfig(req *AddWebhookConfigRequest) WebHookConfigRead {
	if req.Name == nil {
		req.Name = utils.Adr(uuid.New().String()[:6])
	}
	if req.Active == nil {
		req.Active = utils.Adr(true)
	}
	if req.When != nil {
		globaldata.CheckCondition(*req.When)
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
func DalAddWebhookConfig(req *AddWebhookConfigRequest) (*model.WebHook, *model.WebHookCondition, []model.NeedNode) {
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
				condition, err := query.Q.WebHookCondition.Where(query.Q.WebHookCondition.ID.Eq(*req.ConditionId)).First()
				if err != nil || condition == nil {
					return core.NewKnownError(http.StatusBadRequest, err, "condition not exist")
				}
				webhook.WebHookConditionRefer = req.ConditionId
			} else {
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
// @Router /api/v1/webhook/{id} [get]
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
	out.NeedNodeList = GetNeedeNodesListByWebhookId(id)
	return out
}

func GetNeedeNodesListByWebhookId(id int64) []string {
	var needNodeList []string
	needNodes := DalGetNeedNodesByWebhookId(id)
	for _, needNode := range needNodes {
		needNodeList = append(needNodeList, needNode.NodeName)
	}
	return needNodeList
}

// GetWebhookConfig router -

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

func DalGetNeedNodesByWebhookId(id int64) []*model.NeedNode {
	var needNodes []*model.NeedNode
	q := query.Q.NeedNode
	needNodes, err := q.Where(q.WebHookRefer.Eq(id)).Find()
	if err != nil {
		slog.Error(utils.WrapError(err).Error())
		panic(core.NewKnownError(core.EntityNotFound, err, "need node not found"))
	}
	return needNodes
}

// GetWebhookConfigByName router ------------------------------
// @Summary 根据名称获取webhook配置
// @Description 根据名称获取webhook配置
// @Tags Webhook
// @Accept  json
// @Produce  json
// @Param name query GetWebhookConfigByNameRequest true "webhook名称"
// @Success 200 {object} GetWebhookConfigByNameResponse
// @Router /api/v1/webhook [get]
func GetWebhookConfigByName(c *gin.Context) {
	var req GetWebhookConfigByNameRequest

	// 入参校验
	core.BindParamAndValidate(c, &req)

	// 逻辑处理
	webhook := DalGetWebhookConfigByName(req.Name)

	// 出参序列化以及校验
	out := core.SerializeData(webhook, &WebHookConfigRead{})
	out.NeedNodeList = GetNeedeNodesListByWebhookId(webhook.ID)
	core.ValidateSchema(out)

	core.SuccessHandler(c, GetWebhookConfigByNameResponse{
		Code:    200,
		Data:    out,
		Message: "Webhook configuration get successfully",
	})
}

type GetWebhookConfigByNameRequest struct {
	Name string `json:"name" form:"name" example:"webhook1"` // webhook名称
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
// @Description                     "def"
// @Description                 ]
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
// @Description ## 参数示例3 : 一直触发
// @Description ```json
// @Description {
// @Description     "rule": {
// @Description         "type": "all-time"
// @Description     }
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
	// check rule
	{
		var condition globaldata.Condition
		utils.DeserializeData(req, &condition)
		globaldata.CheckCondition(condition)
	}
	// to db
	var condition = DalCreateCondition(req)
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

// GetAllConditionsByPage router -
// @Summary 获取触发条件列表
// @Description # 参数说明
// @Description ## 请求参数
// @Description | 参数名称 | 类型 | 必填 | 描述 |
// @Description | --- | --- | --- | --- |
// @Description | page | int | 是 | 页码 |
// @Description | page_size | int | 是 | 每页数量 |
// @Description ## 返回参数
// @Description | 参数名称 | 类型 | 描述 |
// @Description | --- | --- | --- |
// @Description | id | int | 条件ID |
// @Description | condition | string | 条件表达式 |
// @Description | created_at | time.Time | 创建时间 |
// @Description | updated_at | time.Time | 更新时间 |
// @Tags Webhook
// @Accept  json
// @Produce  json
// @Param body body GetAllConditionsByPageRequest true "获取条件列表"
// @Success 200 {object} GetAllConditionsByPageResponse
// @Router /api/v1/webhook/conditions [get]
func GetAllConditionsByPage(c *gin.Context) {
	// 入参校验
	var req GetAllConditionsByPageRequest
	core.BindParamAndValidate(c, &req)

	// 逻辑处理
	conditions, total := ServiceGetAllConditionsByPage(req)

	// 出参序列化以及校验
	var webhookData []WebHookConditionRead
	for _, condition := range conditions {
		webhookData = append(webhookData, core.SerializeData(condition, &WebHookConditionRead{}))
	}

	var outData GetAllConditionsByPageData
	outData.total = total
	outData.conditions = webhookData
	core.ValidateSchema(outData)

	core.SuccessHandler(c, GetAllConditionsByPageResponse{
		Code:    200,
		Data:    outData,
		Message: "Condition get successfully",
	})
}

type GetAllConditionsByPageRequest struct {
	Page     int `json:"page" form:"page" validate:"required"`
	PageSize int `json:"page_size" form:"page_size" validate:"required"`
}

type GetAllConditionsByPageResponse struct {
	Code    int                        `json:"code" example:"200"`
	Data    GetAllConditionsByPageData `json:"data" `
	Total   int                        `json:"total" `
	Message string                     `json:"message" example:"Condition get successfully"`
}

type GetAllConditionsByPageData struct {
	conditions []WebHookConditionRead `json:"conditions" form:"conditions" validate:"required"`
	total      int64                  `json:"total" form:"total" validate:"required"`
}

func ServiceGetAllConditionsByPage(req GetAllConditionsByPageRequest) ([]WebHookConditionRead, int64) {
	var conditions []WebHookConditionRead
	var total int64
	var err error
	// 逻辑处理
	var conditionsDB []*model.WebHookCondition
	conditionsDB, total, err = DalGetAllConditionsByPage(req.Page, req.PageSize)
	if err != nil {
		slog.Error(utils.WrapError(err).Error())
		panic(core.NewKnownError(core.EntityNotFound, err, "conditions not found"))
	}
	// 出参序列化以及校验
	for _, condition := range conditionsDB {
		out := core.SerializeData(condition, &WebHookConditionRead{})
		conditions = append(conditions, out)
	}
	return conditions, total
}

func DalGetAllConditionsByPage(page int, pageSize int) ([]*model.WebHookCondition, int64, error) {
	var conditions []*model.WebHookCondition
	var count int64
	var err error
	q := query.Q.WebHookCondition
	conditions, count, err = q.FindByPage((page-1)*pageSize, pageSize)
	if err != nil {
		slog.Error(utils.WrapError(err).Error())
		return nil, 0, err
	}
	return conditions, count, nil
}

func GetConditionById(c *gin.Context) {
	// 入参校验
	id := c.Param("id")
	if id == "" {
		panic(core.NewKnownError(http.StatusBadRequest, nil, "id is empty"))
	}

	out := ServiceGetConditionById(id)
	core.ValidateSchema(out)

	core.SuccessHandler(c, GetConditionByIdResponse{
		Code:    200,
		Data:    out,
		Message: "Condition get successfully",
	})
}

type GetConditionByIdResponse struct {
	Code    int                  `json:"code" example:"200"`
	Data    WebHookConditionRead `json:"data" `
	Message string               `json:"message" example:"Condition get successfully"`
}

func ServiceGetConditionById(id string) WebHookConditionRead {
	// 逻辑处理
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic(core.NewKnownError(http.StatusBadRequest, err, "id is not int"))
	}
	condition := DalGetConditionById(intId)

	// 出参序列化以及校验
	out := core.SerializeData(condition, &WebHookConditionRead{})
	core.ValidateSchema(out)

	return out
}

func DalGetConditionById(id int64) *model.WebHookCondition {
	var condition *model.WebHookCondition
	var err error
	q := query.Q.WebHookCondition
	condition, err = q.Where(q.ID.Eq(id)).First()
	if err != nil {
		slog.Error(utils.WrapError(err).Error())
		panic(core.NewKnownError(core.EntityNotFound, err, "condition not found"))
	}
	return condition
}

// UpdateCondition router -
// @Summary 更新触发条件
// @Description # 参数说明
// @Description ## 请求参数
// @Description | 参数名称 | 类型 | 必填 | 描述 |
// @Description | --- | --- | --- | --- |
// @Description | id | int | 是 | 条件ID |
// @Description | and | list中嵌套本参数 | 否 | 规则列表，逻辑与 |
// @Description | or | list中嵌套本参数 | 否 | 规则列表，逻辑或 |
// @Description | rule | Rule | 否 | 规则 |
// @Description ## Rule类型 定义
// @Description | 字段 | 类型 | 是否必填 | 描述 |
// @Description | --- | --- |
// @Description | node_name | string | 是 | 节点名称 |
// @Description | type | string | 是 | 规则类型，支持eq ne gt lt all-time in not-in |
// @Description | value | any | 是 | 比对值 |
// @Description ## 参数示例1 : 更新条件ID为1的条件，将节点MyVariable大于123改为小于123
// @Description ```json
// @Description {
// @Description     "id": 1,
// @Description     "rule": {
// @Description         "node_name": "MyVariable",
// @Description         "type": "lt",
// @Description         "value": 123
// @Description     }
// @Description }
// @Description ```
// @Description ## 参数示例2 : 更新条件ID为2的条件，将节点node1等于在["abc","def"]，并且节点node2等于123改为节点node1等于在["abc","def"]，并且节点node2等于123
// @Description ```json
// @Description {
// @Description     "id": 2,
// @Description     "and": [
// @Description         {
// @Description             "rule": {
// @Description                 "node_name": "node1",
// @Description                 "type": "in",
// @Description                 "value": [
// @Description                     "abc",
// @Description                     "def"
// @Description                 ]
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
// @Description ## 参数示例3 : 更新条件ID为3的条件，将一直触发改为节点node1等于123
// @Description ```json
// @Description {
// @Description     "id": 3,
// @Description     "rule": {
// @Description         "node_name": "node1",
// @Description         "type": "eq",
// @Description         "value": 123
// @Description     }
// @Description }
// @Description ```
// @Description *注意：Condition是嵌套类型，Condition包含and，or，rule，所以and里面可以嵌套and。。。无限嵌套*
// @Tags Webhook
// @Accept  json
// @Produce  json
// @Param id path int true "条件ID"
// @Param body body UpdateConditionRequest true "更新条件"
// @Success 200 {object} UpdateConditionResponse
// @Router /api/v1/webhook/condition/{id} [put]
func UpdateCondition(c *gin.Context) {
	// 入参校验
	id := c.Param("id")
	if id == "" {
		panic(core.NewKnownError(http.StatusBadRequest, nil, "id is empty"))
	}
	var req UpdateConditionRequest
	core.BindParamAndValidate(c, &req)

	// 逻辑处理
	condition := ServiceUpdateCondition(id, &req)

	// 出参序列化以及校验
	out := core.SerializeData(condition, &WebHookConditionRead{})
	core.ValidateSchema(out)

	core.SuccessHandler(c, UpdateConditionResponse{
		Code:    200,
		Data:    out,
		Message: "Condition updated successfully",
	})
}

type UpdateConditionRequest struct {
	And  []globaldata.Condition `json:"and" form:"and"`   // 规则列表，逻辑与
	Or   []globaldata.Condition `json:"or" form:"or"`     // 规则列表，逻辑或
	Rule *globaldata.Rule       `json:"rule" form:"rule"` // 规则
}

type UpdateConditionResponse struct {
	Code    int                  `json:"code" example:"200"`
	Data    WebHookConditionRead `json:"data" `
	Message string               `json:"message" example:"Condition updated successfully"`
}

func ServiceUpdateCondition(id string, req *UpdateConditionRequest) WebHookConditionRead {
	// check params
	if req.And == nil && req.Or == nil && req.Rule == nil {
		panic(core.NewKnownError(http.StatusBadRequest, nil, "params is empty"))
	}
	// check rule
	{
		var condition globaldata.Condition
		utils.DeserializeData(req, &condition)
		globaldata.CheckCondition(condition)
	}
	// to db
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic(core.NewKnownError(http.StatusBadRequest, err, "id is not int"))
	}
	DalUpdateCondition(intId, req)
	var condition = DalGetConditionById(intId)
	out := core.SerializeData(condition, &WebHookConditionRead{}) // orm model -> out
	return out
}

func DalUpdateCondition(id int64, req *UpdateConditionRequest) *model.WebHookCondition {
	var condition = model.WebHookCondition{}
	var err = query.Q.Transaction(func(tx *query.Query) error {
		condition.Condition = utils.PrintDataAsJson(req)
		_, err := tx.WebHookCondition.Select(tx.WebHookCondition.Condition).Where(tx.WebHookCondition.ID.Eq(id)).Updates(condition)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		sqlErr, ok := err.(*pgconn.PgError)
		if ok {
			slog.Error(utils.WrapError(err).Error())
			panic(core.NewKnownError(core.FieldNotUnique, err, sqlErr.Message))
		}
		slog.Error(utils.WrapError(err).Error())
		panic(core.NewKnownError(core.FieldNotUnique, nil, err.Error()))
	}

	condition.ID = id

	return &condition
}

// DeleteCondition router -
// @Summary 删除触发条件
// @Tags Webhook
// @Accept  json
// @Produce  json
// @Param id path int true "条件ID"
// @Success 200 {object} DeleteConditionResponse
// @Router /api/v1/webhook/condition/{id} [delete]
func DeleteCondition(c *gin.Context) {
	// 入参校验
	id := c.Param("id")
	if id == "" {
		panic(core.NewKnownError(http.StatusBadRequest, nil, "id is empty"))
	}

	ServiceDeleteCondition(id)

	core.SuccessHandler(c, DeleteConditionResponse{
		Code:    200,
		Message: "Condition deleted successfully",
	})
}

type DeleteConditionResponse struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"Condition deleted successfully"`
}

func ServiceDeleteCondition(id string) {
	// 逻辑处理
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic(core.NewKnownError(http.StatusBadRequest, err, "id is not int"))
	}
	DalDeleteCondition(intId)
}

func DalDeleteCondition(id int64) {
	err := query.Q.Transaction(func(tx *query.Query) error {
		var err error
		_, err = tx.WebHookCondition.Where(tx.WebHookCondition.ID.Eq(id)).Delete()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		sqlErr, ok := err.(*pgconn.PgError)
		if ok {
			slog.Error(utils.WrapError(err).Error())
			panic(core.NewKnownError(core.FieldNotUnique, err, sqlErr.Message))
		}
		slog.Error(utils.WrapError(err).Error())
		panic(core.NewKnownError(core.FieldNotUnique, err, sqlErr.Message))
	}
}
