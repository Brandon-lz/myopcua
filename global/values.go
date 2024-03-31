package globaldata

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/Brandon-lz/myopcua/db/gen/model"
	"github.com/Brandon-lz/myopcua/db/gen/query"
	"github.com/Brandon-lz/myopcua/utils"
)

var (
	Debug       bool = false
	OPCNodeVars *OPCNodeVarsDFT
	WebHooks    *WebHookConditions
)

func InitSystemVars() {
	// 判断systemvars.obj是否存在，不存在则创建
	if _, err := os.Stat("systemvars.obj"); os.IsNotExist(err) {
		OPCNodeVars = NewGlobalOPCNodeVars()
		if err := OPCNodeVars.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "save systemvars failed:%v\n", err)
		}

	} else {
		if err := utils.LoadObj("systemvars.obj", &OPCNodeVars); err != nil {
			log.Fatalf("load systemvars failed:%v\n", err)
		}
	}

	// todo: 初始化WebHooks
	initWebHooks()
}

func initWebHooks() {
	WebHooks = NewWebHookConditions()

	var webhooks []*model.WebHook
	var err error
	webhooks, err = query.Q.WebHook.Find()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	for _, hook := range webhooks {
		if hook.WebHookConditionRefer == nil || !*hook.Active {
			continue
		}
		wh := utils.DeserializeData(hook, &WebHookConfig{})
		dbConditions, err := query.Q.WebHookCondition.Where(query.WebHookCondition.ID.Eq(*hook.WebHookConditionRefer)).Find()
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		if len(dbConditions) == 0 {
			continue
		}
		wh.ConditionId = hook.WebHookConditionRefer
		utils.DeserializeData(dbConditions[0].Condition, &wh.When)

		needNodes, err := query.Q.NeedNode.Where(query.NeedNode.WebHookRefer.Eq(hook.ID)).Find()
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		for _, node := range needNodes {
			wh.NeedNodeNameList = append(wh.NeedNodeNameList, node.NodeName)
		}

		slog.Debug(fmt.Sprintf("111111111111111111add webhook:%+v", wh.Active))
		WebHooks.AddWebHookConfig(&wh)
	}
}
