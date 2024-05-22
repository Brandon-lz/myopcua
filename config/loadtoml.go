package config

import (
	"log"
	"os"

	"github.com/Brandon-lz/myopcua/utils"
	"github.com/BurntSushi/toml"
)

func Init(configfile string) {
	// Register the config loader
	if err := loadConfig(configfile); err != nil {
		log.Fatalln("Failed to load config file: ", err)
		os.Exit(1)
	}
}

func loadConfig(tomlFilePath string) error {
	// Load the config file
	var configmap map[string]interface{}
	if _, err := toml.DecodeFile(tomlFilePath, &configmap); err != nil {
		// 处理错误
		log.Fatalln("Failed to load config file: ", err)
		return err
	}

	utils.DeserializeData(configmap,&Config)

	if os.Getenv("run_env") != "" {
		Config.RunEnv = os.Getenv("run_env")
	}

	return nil
}
