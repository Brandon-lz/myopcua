package test

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"net/url"

	"github.com/stretchr/testify/require"
)

func PostRequest(urlstr string, body string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlstr, bytes.NewBuffer([]byte(body)))
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("accept", "application/json")
	return client.Do(req)
}

type QueryParam struct {
	Name  string
	Value string
}

func GetRequest(urlstr string,query ...QueryParam) (*http.Response, error) {
	client := &http.Client{}

	for i,q := range query{
		if i==0{
			urlstr+="?"+url.QueryEscape(q.Name)+"="+url.QueryEscape(q.Value)
		}else{
			urlstr+="&"+url.QueryEscape(q.Name)+"="+url.QueryEscape(q.Value)
		}
	}
	req, _ := http.NewRequest("GET", urlstr, nil)
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("accept", "application/json")
	return client.Do(req)
}

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
				"node_name": "node1",
				"type": "gt",
				"value": "123" 
			}
		}
	}
	`
	res, err := PostRequest(url, body)
	require.NoError(err)
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("accept", "application/json")

	require.Equal(http.StatusOK, res.StatusCode) // 断言状态码

	resData, err := io.ReadAll(res.Body)
	require.NoError(err)

	t.Log("resoponseData:" + string(resData))
}

func testGetWebhookConfigById(t *testing.T) {
	require := require.New(t)
	// web,condition:=webhookrouters.DalGetWebhookConfigById(1)
	// require.NotNil(web)
	// require.NotNil(condition)
	// t.Log(web)
	res, err := GetRequest("http://localhost:8080/api/v1/webhook/config/1")
	require.NoError(err)
	require.Equal(http.StatusOK, res.StatusCode) // 断言状态码
}
