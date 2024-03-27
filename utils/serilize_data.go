package utils

import (
	"encoding/json"
	"github.com/Brandon-lz/myopcua/log"
)

func SerializeData[T interface{}](source any, target *T) T { // target必须为指针类型
	jsonData, err := json.Marshal(source)
	if err != nil {
		log.Logger.Error("JSON序列化失败: %s", WrapError(err))
		panic(err)
	}

	err = json.Unmarshal(jsonData, target)
	if err != nil {
		log.Logger.Error("JSON反序列化失败: %s", WrapError(err))
		panic(err)
	}
	return *target
}

