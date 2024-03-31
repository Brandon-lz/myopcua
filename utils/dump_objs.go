package utils

import (
	"encoding/gob"
	"fmt"
	"os"
	"reflect"
)

func DumpObj2Local(obj interface{}, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("创建文件失败：%s", err.Error())
	}
	defer file.Close()

	enc := gob.NewEncoder(file)

	err = enc.Encode(obj)
	if err != nil {
		return fmt.Errorf("序列化失败： %s", err.Error())
	}
	return nil
}

func LoadObj(filepath string, obj interface{}) error {
	if reflect.ValueOf(obj).Kind() != reflect.Ptr || reflect.ValueOf(obj).IsNil() {
		return fmt.Errorf("parameter obj must be a non-nil pointer")
	}

	file2, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("打开文件失败： %s", err.Error())
	}
	defer file2.Close()

	dec := gob.NewDecoder(file2)

	err = dec.Decode(obj)
	if err != nil {
		return fmt.Errorf("反序列化失败： %s", err.Error())
	}
	return nil
}
