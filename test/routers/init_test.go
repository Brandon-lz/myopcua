package test

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/Brandon-lz/myopcua/config"
	sysdb "github.com/Brandon-lz/myopcua/db"
	gentool "github.com/Brandon-lz/myopcua/db/gen"
	"github.com/Brandon-lz/myopcua/db/gen/query"
	globaldata "github.com/Brandon-lz/myopcua/global"
	httpservice "github.com/Brandon-lz/myopcua/http_service"
	"github.com/Brandon-lz/myopcua/log"
	opcservice "github.com/Brandon-lz/myopcua/opc_service"
	"github.com/stretchr/testify/assert"
)


func TestMain(t *testing.T) {
	assert := assert.New(t)
	assert.NoError(os.Remove("./systemvars.obj"))
  	
	sysdb.InitDB()
	cleanDb(assert)
	migrateModel()

	gentool.RunGen("../../db/gen/query")

	config.Init("./config.toml")
	log.Init(slog.LevelDebug)
	slog.Info("Starting the opc application...")
	sysdb.InitDB()
	query.SetDefault(sysdb.DB) // init gen model, for decouple with db
	globaldata.InitSystemVars()
	go opcservice.Start()
	go httpservice.Start()

	time.Sleep(1 * time.Second)

	// Run router tests ---------------------------------
	// t.Run("Test_AddWebhookConfig", testAddWebhookConfig)
	// t.Run("Test_getWebhookConfig", testGetWebhookConfigById)
	testAddWebhookConfig(t)
	testGetWebhookConfigById(t)
}

func cleanDb(assert *assert.Assertions) {
	res := sysdb.DB.Exec("DROP SCHEMA public CASCADE;CREATE SCHEMA public;")
	assert.NoError(res.Error)
}

func migrateModel() {
	sysdb.InitDB()
}
