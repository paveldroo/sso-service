package main

import (
	"log/slog"

	"github.com/paveldroo/sso-service/internal/config"
	"github.com/paveldroo/sso-service/internal/db/sqlite"
	"github.com/paveldroo/sso-service/internal/logger/sl"
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

	// TODO: logger

	// TODO: storage

	// TODO: handlers

	// TODO: build grpc app

	// TODO: grpc app
}
