package log

import (
	"fmt"

	"github.com/Brandon-lz/myopcua/utils"
)

type ElasticsearchWriter struct {
}

func NewElasticsearchWriter() *ElasticsearchWriter {
	return &ElasticsearchWriter{
	}
}

func (w *ElasticsearchWriter) Write(p []byte) (n int, err error) {
	go func() {
		defer utils.RecoverAndLog()
		fmt.Println(string(p))
	}()
	return len(p), nil
}
