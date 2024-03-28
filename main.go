package main

import (
	"log/slog"

	"github.com/Brandon-lz/myopcua/config"
	"github.com/Brandon-lz/myopcua/db"
	"github.com/Brandon-lz/myopcua/db/gen/query"
	globaldata "github.com/Brandon-lz/myopcua/global_data"

	// "github.com/Brandon-lz/myopcua/health"
	"github.com/Brandon-lz/myopcua/log"

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
	config.Init("./config.toml")
	log.Init(slog.LevelDebug)
	slog.Info("Starting the opc application...")
	globaldata.InitSystemVars()
	go opcservices.Start()
	db.InitDB()
	query.SetDefault(db.DB)    // init gen model, for decouple with db
	go httpservice.Start()
	select {}
}
