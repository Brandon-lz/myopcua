package globaldata

import (
	"fmt"
	"log"
	"os"

	"github.com/Brandon-lz/myopcua/utils"
)

var (
	OPCNodeVars *OPCNodeVarsDFT
	WebHook     *WebHookConfig
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
}
