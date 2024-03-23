package test

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTry(t *testing.T) {
	assert := assert.New(t)
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
