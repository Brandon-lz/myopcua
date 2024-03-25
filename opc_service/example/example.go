package main

import (
	"context"
	"fmt"
	"log"

	"time"

	opcUa "github.com/Brandon-lz/myopcua/opc_ua"

	"github.com/gopcua/opcua/ua"
)

func main() {
	ctx := context.Background()
	c, err := opcUa.ConnectToDevice(ctx, "opc.tcp://opcserver:4840", false)
	if err != nil {
		panic(err)
	}
	defer c.Close(ctx)

	nodeId := "ns=2;i=2"

	value, err := opcUa.ReadValueByNodeId(nodeId, ctx, c)
	if err != nil {
		panic(err)
	}
	log.Println("read value:", value)

	values, err := opcUa.ReadMultiValueByNodeIds([]string{nodeId, nodeId}, []*ua.ReadValueID{}, ctx, c)
	if err != nil {
		panic(err)
	}
	log.Println("read values:", values)

	node, err := opcUa.GetNodeByNodeId(nodeId, c)
	if err != nil {
		panic(err)
	}

	rootnode, err := opcUa.GetRootNode(c)
	if err != nil {
		panic(err)
	}
	log.Println("root node:", rootnode)

	children, err := opcUa.GetNodeChildren(rootnode, ctx)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("get node children", children)
	for _, v := range children {
		name, _ := opcUa.GetNodeDisplayName(v, ctx)
		fmt.Println(name)
		// cc, _ := opcUa.GetNodeChildren(v, ctx)
		// for _, vi := range cc {
		// 	name, _ = opcUa.GetNodeDisplayName(vi, ctx)
		// 	fmt.Println("-------------")
		// 	fmt.Println(name)
		// 	cc, _ := opcUa.GetNodeChildren(vi, ctx)
		// 	for _, vii := range cc {
		// 		name, _ = opcUa.GetNodeDisplayName(vii, ctx)
		// 		fmt.Println(name)
		// 	}
		// }
	}

	name, err := opcUa.GetNodeDisplayName(node, ctx)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("节点名称", name)

	description, err := opcUa.GetNodeDescripe(node, ctx)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("节点描述", description)

	log.Println("节点类型", node.ID.Type().String())

	uavar, err := node.Value(ctx)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("value from node", uavar.Value())
	}
}

func startOneThread(c chan interface{}, name int) {
	for {
		select {
		case i := <-c:
			switch t := i.(type) {
			case int:
				fmt.Println("got a int ", t)
			case string:
				fmt.Println("got a string", t)
			case map[string]interface{}:
				fmt.Println("got a map", t)
			}
			// fmt.Println(i)
			// return
			return
		default:
			fmt.Println(name)
			time.Sleep(200 * time.Millisecond)
			// return
		}
	}
}
