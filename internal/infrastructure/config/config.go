package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

const (
	StorageTypeMemory   string = "memory"
	StorageTypePostgres string = "postgres"
)

type Config struct {
	Port        string `env:"PORT" envDefault:"8080"`
	StorageType string `env:"STORAGE_TYPE" envDefault:"memory"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	var storageFlag string
	flag.StringVar(&storageFlag, "storage", "", "storage type: memory or postgres")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	if storageFlag != "" {
		cfg.StorageType = storageFlag
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	switch c.StorageType {
	case StorageTypeMemory, StorageTypePostgres:
		return nil
	default:
		return fmt.Errorf("invalid storage type: must be 'memory' or 'postgres', got %s", c.StorageType)
	}
}
