package config

type GlobalConfig struct {
	Opcua   Opcua   `json:"opcua"`
	Openapi Openapi `json:"openapi"`
}

type Opcua struct {
	Endpoint string `json:"endpoint"`
}

type Openapi struct {
	Port int `json:"port"`
}
