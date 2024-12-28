package serverConfig

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type ServerConfig struct {
	Env        string `yaml:"env" env:"ENV" env-default:"local"`
	HttpServer `yaml:"http_server" env-required:"true"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	TimeOut     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoadServerConfig() *ServerConfig {

	configPath := os.Getenv("CONFIG_SERVER_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_SERVER_PATH environment variable not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_SERVER_PATH does not exist %s", configPath)
	}

	var config ServerConfig

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("Cannot load config file: %s", err)
	}

	return &config
}
