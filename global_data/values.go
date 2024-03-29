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
		if err := utils.Load("systemvars.obj", &OPCNodeVars); err != nil {
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
	webhooks,err = query.Q.WebHook.Find()
	if err!= nil {
		slog.Error(err.Error())
		panic(err)
	}

	for _, hook := range webhooks {
		wh := utils.SerializeData(hook, &WebHookConfig{})
		WebHooks.AddWebHookConfig(&wh)
	}
}

