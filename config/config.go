package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	Env          string `env:"ENV" envDefault:"dev"`
	Port         int    `env:"PORT" envDefault:"8080"`
	KVSHost      string `env:"KVS_HOST" envDefault:"127.0.0.1"`
	KVSPort      int    `env:"KVS_PORT" envDefault:"6379"`
	KVSExpireSec int    `env:"KVS_EXPIRE_SEC" envDefault:"86400"`
	TOKENSECRETS string `env:"TOKEN_SECRETS"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
