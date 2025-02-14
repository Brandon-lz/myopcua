package main

import (
	"context"
	"fmt"
	"log"

	"time"

	opcUa "github.com/Brandon-lz/myopcua/opc_ua"
)

func main() {
	ctx := context.Background()
	c, err := opcUa.ConnectToDevice(ctx, "opc.tcp://192.168.70.173:34840", false)
	if err != nil {
		panic(err)
	}
	defer c.Close(ctx)

	rootnode, err := opcUa.GetRootNode(c)
	if err != nil {
		panic(err)
	}
	log.Println("root node:", rootnode)

	plcnode, err := opcUa.GetNodeByNodeId("ns=3;s=ServerInterfaces", c)
	if err != nil {
		panic(err)
	}

	res := opcUa.SearchChildren("aim", plcnode, ctx)
	fmt.Println(res)
	for _, n := range res {
		name, _ := opcUa.GetNodeDisplayName(&n, ctx)
		fmt.Println(name)
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
