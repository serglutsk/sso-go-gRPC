package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	GRPC        GRPCConfig    `yaml:"grpc"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTtl    time.Duration `yaml:"token_ttl" env-required:"true"`
}

type GRPCConfig struct {
	Timeout time.Duration `yaml:"timeout"`
	Port    int           `yaml:"port"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file doesn't exist: " + path)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("Failed to read config: " + err.Error())
	}
	return &cfg
}

func fetchConfigPath() string {
	var result string
	//--config="path/to/config.yaml"
	flag.StringVar(&result, "config", "", "path to config file")
	flag.Parse()
	if result == "" {
		result = os.Getenv("CONFIG_PATH")
	}
	return result
}
