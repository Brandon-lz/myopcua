package config

type GlobalConfig struct {
	RunEnv  string  `json:"run_env"`
	Opcua   Opcua   `json:"opcua"`
	Openapi Openapi `json:"openapi"`
}

type Opcua struct {
	Endpoint string `json:"endpoint"`
}

type Openapi struct {
	Port int `json:"port"`
}
