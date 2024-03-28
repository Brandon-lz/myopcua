package utils

import (
	"encoding/json"
	"log/slog"
)

func SerializeData[T interface{}](source any, target *T) T { // target必须为指针类型
	var jsonData []byte
	var err error

	if sourceString, isString := source.(string); isString {
		jsonData = []byte(sourceString)
	} else if sourceBytes, isBytes := source.([]byte); isBytes {
		jsonData = sourceBytes
	} else {
		jsonData, err = json.Marshal(source)
		if err != nil {
			slog.Error("JSON序列化失败: " + WrapError(err).Error())
			panic(err)
		}
	}

	err = json.Unmarshal(jsonData, target)
	if err != nil {
		slog.Error("JSON反序列化失败: " + WrapError(err).Error())
		panic(err)
	}
	return *target
}
