package config

type GlobalConfig struct {
	RunEnv  string  `json:"run_env"`
	Opcua   Opcua   `json:"opcua"`
	Openapi Openapi `json:"openapi"`
}

type Opcua struct {
	Endpoint string `json:"endpoint"`
	Interval int `json:"interval"` // in milliseconds, default is 100 milliseconds
}

type Openapi struct {
	DeployHost string `json:"deploy_host"`
	Port int `json:"port"`
}
