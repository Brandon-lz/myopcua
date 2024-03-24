package utils

import (
	"encoding/gob"
	"fmt"
	"os"
	"reflect"
)


func Dump(obj interface{},filepath string) error {

    // 打开一个文件用于写入序列化后的数据
    file, err := os.Create(filepath)
    if err != nil {
        return fmt.Errorf("创建文件失败：%s", err.Error())
    }
    defer file.Close()

    // 创建一个编码器
    enc := gob.NewEncoder(file)

    // 将Person对象序列化并写入文件
    err = enc.Encode(obj)
    if err != nil {
        return fmt.Errorf("序列化失败： %s", err.Error())
    }
	return nil

}

func Load(filepath string, obj interface{}) error {
	if reflect.ValueOf(obj).Kind() != reflect.Ptr || reflect.ValueOf(obj).IsNil() {
		return fmt.Errorf("parameter obj must be a non-nil pointer")
	}
	
    // 打开一个文件用于读取序列化后的数据
    file2, err := os.Open(filepath)
    if err != nil {
        return fmt.Errorf("打开文件失败： %s", err.Error())
    }
    defer file2.Close()

    // 创建一个解码器
    dec := gob.NewDecoder(file2)

    // 从文件中读取序列化后的数据并反序列化为Person对象
    err = dec.Decode(obj)
    if err != nil {
        return fmt.Errorf("反序列化失败： %s", err.Error())
    }
	return nil
}
