package config

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/paveldroo/sso-service/internal/logger/sl"
	"gopkg.in/yaml.v3"
)

type Config struct {
	GRPC GRPCConfig `yaml:"grpc"`
	DB   DB         `yaml:"db"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type DB struct {
	Path           string `yaml:"db_path"`
	MigrationsPath string `yaml:"migrations_path"`
}

func MustLoad() *Config {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "", "path to app config file")
	flag.Parse()

	if cfgPath == "" {
		slog.Error("config file path is empty, usage: --config=<path_to_file>")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		slog.Error("config file doesn't exists: " + cfgPath)
		os.Exit(1)
	}

	return MustLoadPath(cfgPath)
}

func MustLoadPath(cfgPath string) *Config {
	buf, err := os.ReadFile(cfgPath)
	if err != nil {
		slog.Error("failed to open config file", sl.Err(err))
		os.Exit(1)
	}

	cfg := &Config{}
	if err = yaml.Unmarshal(buf, cfg); err != nil {
		slog.Error("failed to parse config file", sl.Err(err))
	}

	return cfg
}
