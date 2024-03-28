package globaldata

import (
	"encoding/json"
	"log/slog"

	"github.com/Brandon-lz/myopcua/db/gen/model"
	"github.com/Brandon-lz/myopcua/db/gen/query"
	"github.com/Brandon-lz/myopcua/utils"
)


type Condition struct {
	And  []Condition `json:"and" form:"and"`   // 规则列表，逻辑与
	Or   []Condition `json:"or" form:"or"`     // 规则列表，逻辑或
	Rule *Rule       `json:"rule" form:"rule"` // 规则
}

type Rule struct {
	Type     string  `json:"type" form:"type" binding:"required,oneof=eq ne gt lt all-time in not-in" example:"eq"` // 规则类型 eq, ne, gt, lt, all-time, in, not-in: 相等, 不相等, 大于, 小于, 全时间, 包含, 不包含
	NodeName string  `json:"node_name" form:"node_name" binding:"required" example:"MyVariable"`          // 节点名称
	Value    interface{} `json:"value" form:"value"`                                            // 规则value
}


type WebHookConfig struct {
	Id          int64      `json:"id" form:"id" validate:"required"`
	Name        string     `json:"name" form:"name" validate:"required"`
	Url         string     `json:"url" form:"url" validate:"required"`
	Active      bool       `json:"active" form:"active" validate:"required"`
	When        *Condition `json:"when" form:"when" validate:"omitempty"`
	ConditionId *int64     `json:"condition_id" form:"condition_id" validate:"omitempty"`
}

func GetAllWebHookConfig() ([]*WebHookConfig ,error){
	var out []*WebHookConfig
	tuples, err := DalGetAllWebhookConfig()
	if err != nil {
		slog.Error("get all webhook config error: " + utils.WrapError(err).Error())
		return nil, err
	}
	for _, tuple := range tuples {
		var webhook WebHookConfig = utils.SerializeData(tuple.Webhook, &WebHookConfig{})
		if tuple.Condition != nil {
			webhook.Id = tuple.Condition.ID
			err := json.Unmarshal([]byte(tuple.Condition.Condition), &webhook.When)
			if err != nil {
				slog.Error("unmarshal condition error: " + err.Error())
				return nil, err
			}
			out = append(out, &webhook)
		}
	}
	return out,nil
}


type tuple struct {
	Webhook   *model.WebHook
	Condition *model.WebHookCondition
}

func DalGetAllWebhookConfig() ([]*tuple,error) {
	var out []*tuple
	q := query.Q.WebHook
	u := query.Q.WebHookCondition
	webhooks, err := q.Find()
	if err != nil {
		slog.Error("get all webhook config error: " + utils.WrapError(err).Error())
		return nil, err
	}
	for _, webhook := range webhooks {
		if webhook.WebHookConditionRefer == nil {
			out = append(out, &tuple{Webhook: webhook, Condition: nil})
			continue
		}
		// gen foreign key condition not support yet, deal with handly
		condition, err := u.Where(u.ID.Eq(*webhook.WebHookConditionRefer)).First()
		if err != nil {
			slog.Error(utils.WrapError(err).Error())
			return nil, err
		}
		out = append(out, &tuple{Webhook: webhook, Condition: condition})
	}
	return out,nil
}