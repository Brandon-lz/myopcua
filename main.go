package main

import (
	"earth/config"
	globaldata "earth/global_data"
	"earth/health"
	httpservice "earth/http_service"
	opcservices "earth/opc_service"
	"log"
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
	go httpservice.Start()
	go health.Runhealthcheck()
	select {}
}
