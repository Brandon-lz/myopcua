package utils

import (
	"encoding/json"
	"fmt"
)

func PrintDataAsJson(m interface{}) string {
	d, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("PrintDataAsJson Error:%+v", err))
	}
	fmt.Println(string(d))
	return string(d)
}
