package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Nats struct {
	ConnString string `yaml:"conn_string"`
}

type ElasticSearch struct {
	Addrs    []string `yaml:"addrs"`
	CaPath   string   `yaml:"ca_path"`
	Username string   `env:"ELASTICSEARCH_USERNAME" env-required:"true"`
	Password string   `env:"ELASTICSEARCH_PASSWORD" env-required:"true"`
}

type Postgres struct {
	ConnString string `yaml:"conn_string"`
}

type Config struct {
	Nats          Nats          `yaml:"nats"`
	ElasticSearch ElasticSearch `yaml:"elasticsearch"`
	Postgres      Postgres      `yaml:"postgres"`
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
