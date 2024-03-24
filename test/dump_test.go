package test

import (
	"encoding/gob"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)
type Person struct {
    Name string
    Age  int
}


func TestDump(t *testing.T) {
	require := require.New(t)
	// 创建一个Person对象
    p := &Person{Name: "张三", Age: 30}
    // 打开一个文件用于写入序列化后的数据
    file, err := os.Create("person.obj")
	require.NoError(err)
    defer file.Close()
    // 创建一个编码器
    enc := gob.NewEncoder(file)

    // 将Person对象序列化并写入文件
    err = enc.Encode(p)
    require.NoError(err)

    t.Log("序列化成功！")
}


func TestLoad(t *testing.T) {
	require := require.New(t)
	// 打开一个文件用于读取序列化后的数据
    file, err := os.Open("person.obj")
	require.NoError(err)
    defer file.Close()
    // 创建一个解码器
    dec := gob.NewDecoder(file)

    // 创建一个Person对象
    var p Person
    // 从文件中反序列化Person对象
    err = dec.Decode(&p)
    require.NoError(err)


    t.Log("反序列化成功！")
    t.Log(p)
}