package opcservice

import (
	globaldata "github.com/Brandon-lz/myopcua/global_data"
	"github.com/Brandon-lz/myopcua/utils"
)

func checkWebhook() {
	for conditionId, condition := range globaldata.WebHooks.ConditionList {
		if CheckCondition(*condition){
			globaldata.WebHooks.FindWebHookByConditionId(int64(conditionId)).SendMsg()
		}
	}
}

func CheckCondition(condition globaldata.Condition) bool {
	var subResultAnd bool = true
	var subResultOr bool = false
	var subResultRule bool = false
	if condition.And != nil {
		for _, subCondition := range condition.And {
			subResultAnd = subResultAnd && CheckCondition(subCondition)
		}
	}
	if condition.Or != nil {
		for _, subCondition := range condition.Or {
			subResultOr = subResultOr || CheckCondition(subCondition)
		}
	}
	if condition.Rule != nil {
		subResultRule = CheckRule(*condition.Rule)
	}
	return subResultAnd && subResultOr && subResultRule
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
		result = false
		for _, v := range rule.Value.([]interface{}) {
			if val == v {
				result = true
				break
			}
		}
	case "not-in":
		result = true
		for _, v := range rule.Value.([]interface{}) {
			if val == v {
				result = false
				break
			}
		}
	}
	return result
}
