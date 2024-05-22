package config

// [opcua]
// endpoint = "opc.tcp://192.168.1.10:4840"             # OPC-UA 服务器地址
// interval = 500      # 读取间隔，单位 milliseconds

// [openapi]
// deploy_host = "localhost:8080"       # 部署到目标位置
// port = 8080       # 打开  http://localhost:8080/docs/index.html#/ 访问 API 文档

// [db]
// type = "sqlite"           # "postgres"
// POSTGRES_HOST = "db"
// POSTGRES_USER = "postgres"
// POSTGRES_DB = "postgres"
// POSTGRES_PASSWORD = "postgres"

type GlobalConfig struct {
	RunEnv  string  `json:"run_env"`
	Opcua   Opcua   `json:"opcua"`
	Openapi Openapi `json:"openapi"`
	DB      DB      `json:"db"`
}

type Opcua struct {
	Endpoint string `json:"endpoint"`
	Interval int `json:"interval"` // in milliseconds, default is 100 milliseconds
}

type Openapi struct {
	DeployHost string `json:"deploy_host"`
	Port int `json:"port"`
}

type DB struct{
	Type string `json:"type"`	
	POSTGRES_HOST string `json:"POSTGRES_HOST"`
	POSTGRES_USER string `json:"POSTGRES_USER"`
	POSTGRES_DB string `json:"POSTGRES_DB"`
	POSTGRES_PASSWORD string `json:"POSTGRES_PASSWORD"`
}