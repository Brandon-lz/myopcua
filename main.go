package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Brandon-lz/myopcua/config"
	"github.com/Brandon-lz/myopcua/db"
	"github.com/Brandon-lz/myopcua/db/gen/query"
	globaldata "github.com/Brandon-lz/myopcua/global"

	// "github.com/Brandon-lz/myopcua/health"
	"github.com/Brandon-lz/myopcua/log"

	httpservice "github.com/Brandon-lz/myopcua/http_service"
	opcservices "github.com/Brandon-lz/myopcua/opc_service"
)

// @title OPC-UA Open API
// @version 1.0
// @description OPC-UA转http协议
// @description 两步完成opcua到http协议的转换(查看下面接口中带**步骤号**字样的接口)
// @contact.email advanced_to@163.com
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.Init("./config.toml")
	log.Init(slog.LevelDebug)
	slog.Info("Starting the opc application...")
	db.InitDB()
	query.SetDefault(db.DB) // init gen model, for decouple with db
	globaldata.InitSystemVars()
	ctx := context.Background()
	ctx,cancle := context.WithCancel(ctx)
	defer cancle()
	go opcservices.Start(ctx)
	go httpservice.Start(ctx)

	quit := make(chan os.Signal,1)
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	cancle()
	slog.Info("Server stopped")
}
