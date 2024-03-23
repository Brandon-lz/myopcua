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


type SystemVarsDFT struct {
	CurrentValues map[int64]*OpcNode  // 0 node1, 1 node2, 2 node3...
	NodeNameSets  map[string]struct{} // set of node names
	NodeIdSets    map[string]struct{} // set of node ids  golang中没有集合  nodeid unique
	NodeIdList    []string            // list of node ids
	NodeNameIndex map[string]int64    // node name to index in CurrentValues   node name 索引
}

func NewSystemVarsDFT() *SystemVarsDFT {
	return &SystemVarsDFT{
		CurrentValues: make(map[int64]*OpcNode),
		NodeNameSets:  make(map[string]struct{}),
		NodeIdSets:    make(map[string]struct{}),
		NodeIdList:    make([]string, 0),
		NodeNameIndex: make(map[string]int64),
	}
}

type OpcNode struct {
	NodeID   string
	Name     string
	DataType string
	Value    interface{}
}



func TestDump(t *testing.T) {
	require := require.New(t)
	// 创建一个Person对象
    // p := &Person{Name: "张三", Age: 30}
	p := NewSystemVarsDFT()
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
    var p SystemVarsDFT
    // 从文件中反序列化Person对象
    err = dec.Decode(&p)
    require.NoError(err)


    t.Log("反序列化成功！")
    t.Log(p)
}