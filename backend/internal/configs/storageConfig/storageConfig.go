package storageConfig

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type StorageConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	Username string `yaml:"username" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"db_name" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env-required:"true"`
}

func MustLoadStorageConfig() StorageConfig {

	configPath := os.Getenv("CONFIG_STORAGE_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_STORAGE_PATH environment variable not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_STORAGE_PATH does not exist: %s", configPath)
	}
	var config StorageConfig

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("cannot read database config: %s", err)
	}

	return config
}
