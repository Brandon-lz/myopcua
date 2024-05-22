package db

import (
	"github.com/Brandon-lz/myopcua/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	switch config.Config.DB.Type {
	case "sqlite":
		DB, err = GetSqliteDB()
	case "postgres":
		DB, err = GetPGDB()
	default:
		panic("unsupported database type")
	}

	// DB, err = GetPGDB()
	// DB, err = GetSqliteDB()

	if err != nil {
		panic("failed to connect database")
	}

	initModels()
	DB.AutoMigrate(modelsToMigrate.modelListToAutoMigrate()...)

	// query.SetDefault(DB)
	sqlDB, err := DB.DB()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
}

// 由于gen在不同环境下生成的model可能不同，暂时不使用sqlite
func GetSqliteDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}

func GetPGDB() (*gorm.DB, error) {
	dsn := "host=" + config.Config.DB.POSTGRES_HOST + " user=" + config.Config.DB.POSTGRES_USER + " password=" + config.Config.DB.POSTGRES_PASSWORD + " dbname=" + config.Config.DB.POSTGRES_DB + " port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}
