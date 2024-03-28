package test

import (
	// "io"
	// "net/http"
	// "net/http/httptest"
	"log/slog"
	"testing"
	"time"

	"github.com/Brandon-lz/myopcua/config"
	sysdb "github.com/Brandon-lz/myopcua/db"
	"github.com/Brandon-lz/myopcua/db/gen/query"
	globaldata "github.com/Brandon-lz/myopcua/global_data"
	httpservice "github.com/Brandon-lz/myopcua/http_service"
	"github.com/Brandon-lz/myopcua/log"
	opcservice "github.com/Brandon-lz/myopcua/opc_service"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(t *testing.T) {
	assert := assert.New(t)
	cleanDb(assert)
	migrateModel()

	config.Init("../../config.toml")
	log.Init()
	// log.Logger.Info("Starting the opc application...")
	slog.Info("Starting the opc application...")
	globaldata.InitSystemVars()
	go opcservice.Start()
	sysdb.InitDB()
	query.SetDefault(sysdb.DB)    // init gen model, for decouple with db
	go httpservice.Start()

	time.Sleep(1 * time.Second)
	// Run tests
	t.Run("Test_GetNode", testHomepageHandler)
}

func cleanDb(assert *assert.Assertions) {
	var err error
	db, err = sysdb.GetPGDB()
	assert.NoError(err)
	res := db.Exec("DROP SCHEMA public CASCADE;CREATE SCHEMA public;")
	assert.NoError(res.Error)
}

func migrateModel() {
	sysdb.InitDB() 
}

