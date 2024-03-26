package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Number interface {
	int | float64
}

type Getter[T Number] interface {
	Get() T
}

type Model[T Number] struct {
	Data []T
}

func (m *Model[T]) Push(item T) {
	m.Data = append(m.Data, item)
}

func (m *Model[T]) Get(i int) T {
	return m.Data[i]
}

func TestGeneric(t *testing.T) {
	assert := assert.New(t)
	// passing int as type parameter
	modelInt := Model[int]{Data: []int{1, 2, 3}}
	fmt.Println(modelInt.Data) // [1 2 3]

	assert.Equal(modelInt.Data, []int{1, 2, 3})

	// passing float64 as type parameter
	modelFloat := Model[float64]{Data: []float64{1.1, 2.2, 0.02}}
	fmt.Println(modelFloat.Data) // [1.1 2.2 0.02]

	modelInt.Push(4)
	fmt.Println(modelInt.Data) // [1 2 3 4]

	itemAtOne := modelFloat.Get(1)
	fmt.Println(itemAtOne) // 2.2
}
