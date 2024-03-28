package test

import (
	"testing"
	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	Age  int
}

func TestSlice(t *testing.T) {
	require:= require.New(t)
	s := []int{1, 2, 3, 4, 5}
	t.Log(s)
	s = append(s[:2],s[3:]...)
	t.Log(s)
	require.Equal(s, []int{1, 2, 4, 5})

	s2 := []*TestStruct{{1}, {2}, {3}, {4}, {5}}
	t.Log(s2)
	// st := s2
	// s2 = append(s2[:2],nil)
	// s2 = append(s2,st[3:]...)
	s2[2] = nil
	t.Log(s2)
	require.Equal(s2, []*TestStruct{{1}, {2}, nil, {4}, {5}})


}
