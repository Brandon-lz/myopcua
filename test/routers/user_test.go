package test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)


func testHomepageHandler(t *testing.T) {
	require := require.New(t)
	// mockResponse := `{"message":"Welcome to the Tech Company listing API with Golang"}`
	// r := httpservice.InitRouter()
	// r.POST("/api/v1/webhook/condition", webhookrouters.AddWebhookConfig)

	client := &http.Client{}

	body := `{
		"active": true,
		"name": "webhook1",
		"url": "http://localhost:8080/api/v1/webhook/example",
		"when": {
			"rule": {
				"node_name": "node1",
				"type": "eq",
				"value": "123" 
			}
		}
	}`
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/webhook/config", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	res,err:=client.Do(req)
	require.NoError(err)

	require.Equal(http.StatusOK, res.StatusCode) // 断言状态码

	resData, err := io.ReadAll(res.Body)
	require.NoError(err)

	t.Log("resoponseData:" + string(resData))
}
