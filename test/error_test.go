package test

import (
	// "errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

const badInput = "abc"

type BadInputError struct {
	input string
}

func (e *BadInputError) Error() string {
	return fmt.Sprintf("bad input: %s", e.input)
}

func validateInput(input string) error {
	if input == badInput {
		return fmt.Errorf("validateInput: %w", &BadInputError{input: input})
	}
	return nil
}

func TestError(t *testing.T) {
	require := require.New(t)
	// input := badInput

	// err := validateInput(input)
	// var badInputErr *BadInputError
	// if errors.As(err, &badInputErr) {
	// 	t.Logf("bad input error occured: %s\n", badInputErr)
	// }
	// require.NotNil(err)

    err := BadInputError{input: "bad input"}
    switchError(&err,t)
    require.NotNil(err)

}


func switchError(err error,t *testing.T) {
	switch err.(type) {
	case *BadInputError:
		t.Log("bad input error")
	default:
		t.Log("unknown error")
	}
}