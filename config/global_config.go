package config


type GlobalConfig struct {
	RunEnv  string  `json:"run_env"`
	Opcua   Opcua   `json:"opcua"`
	Openapi Openapi `json:"openapi"`
	DB      DB      `json:"db"`
}

type Opcua struct {
	Endpoint string `json:"endpoint"`
	Interval int `json:"interval"` // in milliseconds, default is 100 milliseconds
	NodeId string `json:"node_id"`
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