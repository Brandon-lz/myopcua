
## 安装

```
go install github.com/swaggo/swag/cmd/swag@latest
go get .
go get -u github.com/swaggo/swag    # 使用最新版swag，防止错误
```

## 运行

```
go test -v -run TestMigrate test/migrate_test.go    # 迁移数据库
go test -v -run TestGenInit ./test/gen_test.go   # 生成orm代码
swag init
go mod tidy   # 更新依赖
go run .
```
