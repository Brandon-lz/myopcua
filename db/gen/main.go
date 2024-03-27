package main

import (
    "github.com/Brandon-lz/myopcua/db"
    "gorm.io/gen"
    "gorm.io/gen/field"
)

// go test -v -run TestGenInit ./test/gen_test.go 

func main() {
    g := gen.NewGenerator(gen.Config{
        OutPath: "./query",
        Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
        FieldNullable: true,  // generate nullable fields
    })

    gormdb, _ := db.GetPGDB()
    g.UseDB(gormdb) // reuse your gorm db

    // Generate basic type-safe DAO API for struct `model.User` following conventions
    g.ApplyBasic(
        g.GenerateAllTable()...,
    )


    // Generate model for foreign key relation tables
    webhooks := g.GenerateModel("web_hooks")
    webhookConditions := g.GenerateModel("web_hook_conditions", gen.FieldRelate(field.HasMany, "WebHooks", webhooks,
        &field.RelateConfig{
            // RelateSlice: true,
            GORMTag: field.GormTag{"foreignKey": []string{"WebHookConditionRefer"}, "references": []string{"ID"}},
        }),
    )

    g.ApplyBasic(
        // g.GenerateAllTable()...,
        webhooks, webhookConditions,
    )

    // Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
    // g.ApplyInterface(func(Querier) {}, model.User{}, model.Company{})

    // Generate the code
    g.Execute()
}
