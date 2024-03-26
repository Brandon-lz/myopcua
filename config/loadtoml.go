package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

func Init() {
	// Register the config loader
	if err := loadConfig("./config.toml"); err != nil {
		log.Fatalln("Failed to load config file: ", err)
		os.Exit(1)
	}
}

func loadConfig(tomlFilePath string) error {
	// Load the config file
	if _, err := toml.DecodeFile(tomlFilePath, &Config); err != nil {
		// 处理错误
		log.Fatalln("Failed to load config file: ", err)
		return err
	}

	if os.Getenv("run_env") != "" {
		Config.RunEnv = os.Getenv("run_env")
	}

	return nil
}
