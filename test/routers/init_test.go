package test

import (
	// "context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/Brandon-lz/myopcua/config"
	sysdb "github.com/Brandon-lz/myopcua/db"
	gentool "github.com/Brandon-lz/myopcua/db/gen"
	"github.com/Brandon-lz/myopcua/db/gen/query"
	globaldata "github.com/Brandon-lz/myopcua/global"

	// httpservice "github.com/Brandon-lz/myopcua/http_service"
	"github.com/Brandon-lz/myopcua/log"
	// opcservice "github.com/Brandon-lz/myopcua/opc_service"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	assert := assert.New(t)
	os.Remove("./systemvars.obj")
	os.Remove("./sqlite.db")

	sysdb.InitDB()
	cleanDb(assert)
	migrateModel()

	gentool.RunGen("../../db/gen/query")

	config.Init("./config.toml")
	log.Init("DEBUG")
	slog.Info("Starting the opc application...")
	sysdb.InitDB()
	query.SetDefault(sysdb.DB) // init gen model, for decouple with db
	globaldata.InitSystemVars()
	// ctx,cancel	 := context.WithCancel(context.Background())

	// go opcservice.Start(ctx)
	// go httpservice.Start(ctx)

	time.Sleep(1 * time.Second)

	// Run router tests ---------------------------------
	t.Run("Test_AddWebhookConfig", testAddWebhookConfig)
	t.Run("Test_getWebhookConfig", testGetWebhookConfigById)

	// cancel()

}

func cleanDb(assert *assert.Assertions) {
	res := sysdb.DB.Exec("DROP SCHEMA public CASCADE;CREATE SCHEMA public;")
	assert.NoError(res.Error)
}

func migrateModel() {
	sysdb.InitDB()
}
