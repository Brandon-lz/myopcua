package utils

import (
	"bytes"
	"net/http"
	"net/url"
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
