package test

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTry(t *testing.T) {
	assert := assert.New(t)
	res := tryTest()
	assert.Equal(5, res)
}

func tryTest() (res int) {
	defer func(){
		res = 5
	}()
	return 3
}
