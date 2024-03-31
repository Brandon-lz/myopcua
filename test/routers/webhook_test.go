package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/Brandon-lz/myopcua/utils"
	"github.com/stretchr/testify/require"
)

func testAddWebhookConfig(t *testing.T) {
    require := require.New(t)
    url := "http://localhost:8080/api/v1/webhook/config"

    body := `
    {
        "active": true,
        "name": "webhook1",
        "url": "http://localhost:8080/api/v1/webhook/example",
        "when": {
            "rule": {
                "node_name": "MyVariable",
                "type": "gt",
                "value": "123" 
            }
        },
		"need_node_list": ["MyVariable", "MyVariable2"]
    }
    `

	body = `
    {
        "active": true,
        "name": "webhook1",
        "url": "http://localhost:8080/api/v1/webhook/example",
        "when": {
            "rule": {
                "type": "all-time"
            }
        },
		"need_node_list": ["MyVariable"]
    }
    `
    res, err := utils.PostRequest(url, body)
    require.NoError(err)

    require.Equal(http.StatusOK, res.StatusCode) // 断言状态码

    resData, err := io.ReadAll(res.Body)
    require.NoError(err)

    t.Log("resoponseData:" + string(resData))
}

func testGetWebhookConfigById(t *testing.T) {
    require := require.New(t)
    res, err := utils.GetRequest("http://localhost:8080/api/v1/webhook/config/1")
    require.NoError(err)
    require.Equal(http.StatusOK, res.StatusCode) // 断言状态码
}
