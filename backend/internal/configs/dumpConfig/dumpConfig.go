package dumpConfig

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type DockerCfg struct {
	ContainerName string `yaml:"container" env-required:"true"`
	Username      string `yaml:"username" env-required:"true"`
	DbName        string `yaml:"db_name" env-required:"true"`
	Prefix        string `yaml:"prefix" env-required:"true"`
	Dir           string `yaml:"dir" env-required:"true"`
	RestorePrefix string `yaml:"restorePrefix" env-required:"true"`
}

func MustLoadDumpConfig() DockerCfg {
	configPath := os.Getenv("CONFIG_DUMP_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_DUMP_PATH environment variable not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_DUMP_PATH does not exist: %s", configPath)
	}
	var config DockerCfg

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("cannot read database config: %s", err)
	}

	return config
}
