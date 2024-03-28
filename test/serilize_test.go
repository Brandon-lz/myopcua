package test

import (
	// "errors"
	"log/slog"
	"testing"

	"github.com/Brandon-lz/myopcua/http_service/core"
	"github.com/Brandon-lz/myopcua/log"
	"github.com/Brandon-lz/myopcua/utils"
	"github.com/stretchr/testify/require"
)

// GetWebhookConfig router 参数定义，字段描述放在字段后面
type AddWebhookConfigRequest struct {
	Name        interface{}    `json:"name" form:"name" example:"webhook1"`                                            // webhook名称，可以为空
	Url         string     `json:"url" form:"url" binding:"required,url" example:"http://192.168.1.1:8800/notify"` // webhook地址
}


func TestSerialize(t *testing.T) {
	require := require.New(t)
	log.Init(slog.LevelDebug)

	body := `
	{
		"name": 123,
		"url": "http://192.168.1.1:8800/notify"
	}
	`

	out:=utils.SerializeData(body, &AddWebhookConfigRequest{})
	core.ValidateSchema(out)
	
	require.Equal(out.Name, float64(123))

}
