package main

import (
	"log/slog"

	"github.com/paveldroo/sso-service/internal/config"
	"github.com/paveldroo/sso-service/internal/logger/sl"
	"github.com/paveldroo/sso-service/internal/storage/sqlite"
)

func main() {
	cfg, err := config.ParseCfg()
	if err != nil {
		slog.Error("parse config", sl.Err(err))
	}

	err = sqlite.New(cfg)
	if err != nil {
		slog.Error("run migrations", sl.Err(err))

	}

	// TODO: storage

	// TODO: handlers

	// TODO: build grpc app

	// TODO: grpc app
}
