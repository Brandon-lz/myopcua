package db

import (
	"github.com/Brandon-lz/myopcua/db/gen/query"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	DB, err := GetPGDB()
	if err != nil {
		panic("failed to connect database")
	}

	initModels()
	DB.AutoMigrate(modelsToMigrate.modelListToAutoMigrate()...)

	query.SetDefault(DB)
	sqlDB, err := DB.DB()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
}

func GetSqliteDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}

func GetPGDB() (*gorm.DB, error) {
	dsn := "host=vector-pg user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}
