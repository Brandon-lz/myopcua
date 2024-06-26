package globaldata

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Brandon-lz/myopcua/db/gen/model"
	"github.com/Brandon-lz/myopcua/db/gen/query"
	"github.com/Brandon-lz/myopcua/utils"
)

type WebHookConditions struct {
	Webhooks                   map[int64]*WebHookConfig `json:"webhooks"` // 0:webhook1, 1:webhook2, 2:webhook3  注意，这里的1,2,3并不是webhook的id，而是condition的id，或者说切片id
	ConditionList              []*Condition             `json:"conditions"`
	IndexConditionId2WebHookId map[int64]int64          `json:"index_condition_id_2_web_hook_id"` // find webhook by condition id
	// IndexNodeName2WebHookId    map[string]int64         `json:"index_node_name_2_web_hook_id"`    // find webhook by node name
}

type Condition struct {
	And  []Condition `json:"and" form:"and"`   // 规则列表，逻辑与
	Or   []Condition `json:"or" form:"or"`     // 规则列表，逻辑或
	Rule *Rule       `json:"rule" form:"rule"` // 规则
}

type Rule struct {
	Type     string      `json:"type" form:"type" validate:"required" binding:"required,oneof=eq ne gt lt all-time in not-in" example:"eq"` // 规则类型 eq, ne, gt, lt, all-time, in, not-in: 相等, 不相等, 大于, 小于, 全时间, 包含, 不包含
	NodeName *string      `json:"node_name" form:"node_name" example:"MyVariable"`                    // 节点名称
	Value    interface{} `json:"value" form:"value"`                                                                    // 规则value
}

type WebHookConfig struct {
	Id               int64      `json:"id" form:"id" validate:"required"`
	Name             string     `json:"name" form:"name" validate:"required"`
	Url              string     `json:"url" form:"url" validate:"required"`
	Active           bool       `json:"active" form:"active" validate:"required"`
	When             *Condition `json:"when" form:"when" validate:"omitempty"`
	ConditionId      *int64     `json:"condition_id" form:"condition_id" validate:"omitempty"`
	NeedNodeNameList []string   `validate:"required"`
}

func (wc *WebHookConfig) SendMsg() {
	slog.Debug(fmt.Sprintf("send msg to webhook url: %+v", wc.When))
	var values = make(map[string]interface{})
	for _, needNodeName := range wc.NeedNodeNameList {
		val, err := OPCNodeVars.GetValueByName(needNodeName)
		if err != nil {
			val = nil
		}
		values[needNodeName] = val
	}
	slog.Debug(fmt.Sprintf("send msg to webhook url, values: %+v", values))
	utils.PostRequest(wc.Url,
		utils.PrintDataAsJson(map[string]interface{}{
			"values":    values,
			"timestamp": OPCNodeVars.TimeStamp.UnixMilli(),
		}),
	)
	slog.Info(fmt.Sprintf("send msg to webhook url success: %s values: %v", wc.Url,values))
}

func NewWebHookConditions() *WebHookConditions {
	return &WebHookConditions{
		Webhooks:                   make(map[int64]*WebHookConfig),
		ConditionList:              make([]*Condition, 0),
		IndexConditionId2WebHookId: make(map[int64]int64),
		// IndexNodeName2WebHookId:    make(map[string]int64),
	}
}

// 加载webhook配置到内存，用于判断webhook，所以只会添加有效和激活的webhook
func (w *WebHookConditions) AddWebHookConfig(webhook *WebHookConfig) {
	if !webhook.Active || webhook.ConditionId == nil {
		return
	}
	WebHookWriteLock.Lock()
	defer WebHookWriteLock.Unlock()

	var newConditionId int64
	findNilSeat := false
	for i, condition := range w.ConditionList {
		if condition == nil {
			newConditionId = int64(i)
			w.ConditionList[newConditionId] = webhook.When
			w.IndexConditionId2WebHookId[newConditionId] = newConditionId // 这里重复了，应该改成conditionId2WebHookIdinDb
			findNilSeat = true
			w.Webhooks[newConditionId] = webhook
			break
		}
	}
	if !findNilSeat {
		newConditionId = int64(len(w.ConditionList))
		w.ConditionList = append(w.ConditionList, webhook.When)
		w.Webhooks[newConditionId] = webhook
		w.IndexConditionId2WebHookId[newConditionId] = newConditionId
	}
	// w.IndexNodeName2WebHookId[webhook.When.Rule.NodeName] = webhook.Id    // ??? 多个nodename
}

func (w *WebHookConditions) RemoveWebHookConfig(webhookId int64) {
	WebHookWriteLock.Lock()
	defer WebHookWriteLock.Unlock()
	conditionId := w.Webhooks[webhookId].ConditionId
	w.ConditionList[*conditionId] = nil
	delete(w.IndexConditionId2WebHookId, *conditionId)
	// delete(w.IndexNodeName2WebHookId, w.Webhooks[webhookId].When.Rule.NodeName)
	delete(w.Webhooks, webhookId)
}

func (w *WebHookConditions) FindWebHookByConditionId(conditionId int64) *WebHookConfig {
	if webhookId, ok := w.IndexConditionId2WebHookId[conditionId]; ok {
		return w.Webhooks[webhookId]
	}
	return nil
}

// func (w *WebHookConditions) FindWebHookByNodeName(nodeName string) *WebHookConfig {
// 	if webhookId, ok := w.IndexNodeName2WebHookId[nodeName]; ok {
// 		return w.Webhooks[webhookId]
// 	}
// 	return nil
// }

func GetAllWebHookConfig() ([]*WebHookConfig, error) {
	var out []*WebHookConfig
	tuples, err := DalGetAllWebhookConfig()
	if err != nil {
		slog.Error("get all webhook config error: " + utils.WrapError(err).Error())
		return nil, err
	}
	for _, tuple := range tuples {
		var webhook WebHookConfig = utils.DeserializeData(tuple.Webhook, &WebHookConfig{})
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
	return out, nil
}

type tuple struct {
	Webhook   *model.WebHook
	Condition *model.WebHookCondition
}

func DalGetAllWebhookConfig() ([]*tuple, error) {
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
	return out, nil
}
