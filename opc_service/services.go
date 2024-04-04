package opcservice

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	"os"

	"github.com/Brandon-lz/myopcua/config"
	globaldata "github.com/Brandon-lz/myopcua/global"
	opcUa "github.com/Brandon-lz/myopcua/opc_ua"

	"github.com/gopcua/opcua"
)

// 定义全局的conn管道

func Start() {

	slog.Info("启动OPC服务...")

	for {
		err := TestOpc()
		if err != nil {
			slog.Error("OPCUA发生故障或目标设备失去连接，尝试重启服务，请勿关闭服务")
			slog.Error("故障信息:")
			slog.Error(err.Error())
			fmt.Println(err.Error())
			time.Sleep(time.Second * 3)
			continue
		}
		break
	}
}

func TestOpc() (err error) {
	defer func() {
		// except:
		perr := recover()
		if perr != nil {
			err = fmt.Errorf("panic error:[%+v]\n%s", perr, debug.Stack())
		} else {
			slog.Info("设备正常退出")
		}
	}()

	ctx := context.Background()
	c, err := opcUa.ConnectToDevice(ctx, config.Config.Opcua.Endpoint, false)
	if err != nil {
		panic(err)
	}
	defer c.Close(ctx)

	ticker := time.NewTicker(time.Millisecond * time.Duration(config.Config.Opcua.Interval))
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			if IsExpire() {
				os.Exit(0)
			}
			readOpcData(c)
		}
	}()
	select {}
}

func readOpcData(c *opcua.Client) {
	globaldata.OpcWriteLock.Lock() // 加锁 这期间不允许修改SystemVars
	defer globaldata.OpcWriteLock.Unlock()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	datas, err := opcUa.ReadMultiValueByNodeIds(globaldata.OPCNodeVars.NodeIdList, nil, ctx, c)
	if err != nil {
		panic(err)
	}

	slog.Debug("OPC读取数据成功:" + fmt.Sprintf("%+v", datas))
	globaldata.OPCNodeVars.TimeStamp = time.Now()
	// 写入数据到全局变量
	for i, data := range datas {
		slog.Debug("OPC读取数据成功:" + fmt.Sprintf("%+v", globaldata.OPCNodeVars.CurrentValues))
		slog.Debug("OPC读取数据成功:" + fmt.Sprintf("%+v", globaldata.OPCNodeVars.CurrentNodes))
		globaldata.OPCNodeVars.CurrentValues[int64(i)] = data
		globaldata.OPCNodeVars.CurrentNodes[int64(i)].Value = data
	}
	go checkWebhook()

}

func IsExpire() bool {
	fromTime := "2024-03-18T12:00:00Z"
	ft, _ := time.Parse(time.RFC3339, fromTime)
	now := time.Now()
	deadline := ft.Add(time.Hour * 24 * 30)
	if now.After(deadline) {
		slog.Error("授权已过期，请联系平台续费！  wx:advance_to")
		return true
	}
	return false
}
