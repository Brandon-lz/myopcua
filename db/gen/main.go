package main

import (
	"gorm.io/gen"
	"earth/db"
)


func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	
	gormdb, _ := db.GetPGDB()
	g.UseDB(gormdb) // reuse your gorm db

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(
		g.GenerateAllTable()...
	)

	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	// g.ApplyInterface(func(Querier) {}, model.User{}, model.Company{})

	// Generate the code
	g.Execute()
}
