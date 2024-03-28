package test

import (
	"testing"

	"github.com/Brandon-lz/myopcua/utils"
	"github.com/stretchr/testify/assert"
)


func TestGreater(t *testing.T) {
	assert := assert.New(t)
	res, err := utils.GreaterThan2interface(interface{}(int64(10)), interface{}(0))
	assert.NoError(err)
	assert.Equal(true, res)

	res, err = utils.GreaterThan2interface(interface{}("abc"), interface{}("122"))
	assert.Error(err)
	assert.Equal(false, res)

	// assert.Equal()
}