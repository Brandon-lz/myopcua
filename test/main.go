package main

import (
	"fmt"
)

func main() {
	res := tryTest()
	fmt.Println(res)
	var datain interface{} = 12
	dataout,e := datain.(string)
	fmt.Println(dataout,e)

}

func tryTest() (res int) {
	defer func(){
		res = 5
	}()
	return 3
}
