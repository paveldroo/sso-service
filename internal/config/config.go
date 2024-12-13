package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DB DB `yaml:"db"`
}

type DB struct {
	Path           string `yaml:"db_path"`
	MigrationsPath string `yaml:"migrations_path"`
}

func ParseCfg() (*Config, error) {
	buf, err := os.ReadFile("config.yml")
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}

	cfg := &Config{}
	if err = yaml.Unmarshal(buf, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
