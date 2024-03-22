package utils

import (
	"encoding/json"
	"fmt"
)

func PrintMapAsJson(m interface{})string{
	d,err:=json.Marshal(m)
	if err!=nil{
		panic(fmt.Sprintf("PrintMapAsJson Error:%+v",err))
	}
	fmt.Println(string(d))
	return string(d)
}
