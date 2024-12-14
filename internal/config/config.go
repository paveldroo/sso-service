package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server Server `yaml:"server"`
	DB     DB     `yaml:"db"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DB struct {
	Path           string `yaml:"db_path"`
	MigrationsPath string `yaml:"migrations_path"`
}

func ParseCfg(path string) (*Config, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}

	cfg := &Config{}
	if err = yaml.Unmarshal(buf, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
