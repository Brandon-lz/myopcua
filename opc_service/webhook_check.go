package opcservice

import (
	"fmt"
	"log/slog"

	globaldata "github.com/Brandon-lz/myopcua/global_data"
	"github.com/Brandon-lz/myopcua/utils"
)

func checkWebhook() {
	slog.Debug("checkWebhook")
	if globaldata.WebHooks.ConditionList == nil || len(globaldata.WebHooks.ConditionList) == 0 {
		return
	}
	for conditionId, condition := range globaldata.WebHooks.ConditionList {
		slog.Debug(fmt.Sprintf("checkWebhook conditino: %s", utils.PrintDataAsJson(condition)))
		if CheckCondition(*condition) {
			slog.Debug(fmt.Sprintf("checkWebhook conditionId: %d is true", conditionId))
			globaldata.WebHooks.FindWebHookByConditionId(int64(conditionId)).SendMsg()
		}
	}
}

func CheckCondition(condition globaldata.Condition) bool {
	var result bool = false
	if condition.And != nil {
		result = true
		for _, subCondition := range condition.And {
			result = result && CheckCondition(subCondition)
		}
	}
	if condition.Or != nil {
		for _, subCondition := range condition.Or {
			result = result || CheckCondition(subCondition)
		}
	}
	if condition.Rule != nil {
		if condition.And == nil && condition.Or == nil {
			result = CheckRule(*condition.Rule)
		} else {
			result = result && CheckRule(*condition.Rule)
		}
	}
	return result
}

// type Rule struct {
// 	Type     string      `json:"type" form:"type" binding:"required,oneof=eq ne gt lt all-time in not-in" example:"eq"` // 规则类型 eq, ne, gt, lt, all-time, in, not-in: 相等, 不相等, 大于, 小于, 全时间, 包含, 不包含
// 	NodeName string      `json:"node_name" form:"node_name" binding:"required" example:"MyVariable"`                    // 节点名称
// 	Value    interface{} `json:"value" form:"value"`                                                                    // 规则value
// }

func CheckRule(rule globaldata.Rule) bool {
	var result bool = false
	var val, err = globaldata.OPCNodeVars.GetValueByName(rule.NodeName)
	if err != nil {
		return false
	}
	if val == nil {
		slog.Warn(fmt.Sprintf("CheckRule: %s, %s, val is nil", rule.Type, rule.NodeName))
		return false
	}
	slog.Debug(fmt.Sprintf("CheckRule: %s, %s, %v, %v", rule.Type, rule.NodeName, val, rule.Value))

	switch rule.Type {
	case "eq":
		result = val == rule.Value
	case "ne":
		result = val != rule.Value
	case "gt":
		result, err = utils.GreaterThan2interface(val, rule.Value)
		if err != nil {
			return false
		}
	case "lt":
		result, err = utils.GreaterThan2interface(val, rule.Value)
		if err != nil {
			return false
		}
		result = !result
	case "all-time":
		result = true
	case "in":
		for _, v := range rule.Value.([]interface{}) {
			if val == v {
				result = true
				break
			}
		}
		result = false
	case "not-in":
		for _, v := range rule.Value.([]interface{}) {
			if val == v {
				result = false
				break
			}
		}
		result = true
	}
	return result
}
