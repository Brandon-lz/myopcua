package main

import (
	"log"

	"github.com/Brandon-lz/myopcua/config"
	"github.com/Brandon-lz/myopcua/db"
	"github.com/Brandon-lz/myopcua/db/gen/query"
	globaldata "github.com/Brandon-lz/myopcua/global_data"
	"github.com/Brandon-lz/myopcua/health"

	httpservice "github.com/Brandon-lz/myopcua/http_service"
	opcservices "github.com/Brandon-lz/myopcua/opc_service"
)

// @title OPC-UA Open API
// @version 1.0
// @description OPC-UA转http协议
// @contact.email advanced_to@163.com
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.Init()
	log.Println("Starting the opc application...")
	globaldata.InitSystemVars()
	go opcservices.Start()
	db.InitDB()
	query.SetDefault(db.DB)    // 初始化gen模块，放到main函数为了解耦gen和gorm
	go httpservice.Start()
	go health.Runhealthcheck()
	select {}
}
