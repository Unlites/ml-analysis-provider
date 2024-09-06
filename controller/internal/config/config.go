package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Nats struct {
	ConnString string `yaml:"conn_string"`
}

type Server struct {
	Addr string `yaml:"addr"`
}

type Config struct {
	Nats   Nats   `yaml:"nats"`
	Server Server `yaml:"server"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("CONFIG_PATH env is not set")
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}
