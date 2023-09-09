package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	Env     string `json:"env" envDefault:"dev"`
	Port    int    `json:"port" envDefault:"8080"`
	KVSHost string `json:"kvs_host" envDefault:"127.0.0.1"`
	KVSPort int    `json:"kvs_port" envDefault:"6379"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
