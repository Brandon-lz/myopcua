package opcservice

import (
	"fmt"
	"log/slog"

	globaldata "github.com/Brandon-lz/myopcua/global"
	"github.com/Brandon-lz/myopcua/utils"
)

func checkWebhook() {
	slog.Debug("checkWebhook")
	if globaldata.WebHooks.ConditionList == nil || len(globaldata.WebHooks.ConditionList) == 0 {
		return
	}
	for conditionId, condition := range globaldata.WebHooks.ConditionList {
		slog.Debug(fmt.Sprintf("checkWebhook conditino: %s", utils.PrintDataAsJson(condition)))
		if globaldata.CheckCondition(*condition) {
			slog.Debug(fmt.Sprintf("checkWebhook conditionId: %d is true", conditionId))
			globaldata.WebHooks.FindWebHookByConditionId(int64(conditionId)).SendMsg()
		}
	}
}
