package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

var DB *gorm.DB

func InitDB() {
	DB,err:=GetPGDB()
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&WebHook{})
}


func GetSqliteDB() (*gorm.DB ,error) {
	return gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}

func GetPGDB()(*gorm.DB, error) {
	dsn := "host=vector-pg user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}