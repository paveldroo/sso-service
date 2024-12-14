package main

import (
	"log/slog"
	"os"

	"github.com/paveldroo/sso-service/internal/config"
	"github.com/paveldroo/sso-service/internal/logger/sl"
	"github.com/paveldroo/sso-service/internal/storage/sqlite"
)

func main() {
	cfg, err := config.ParseCfg()
	if err != nil {
		slog.Error("parse config", sl.Err(err))
		os.Exit(1)
	}

	storage, err := sqlite.New(cfg)
	if err != nil {
		slog.Error("create sqlite storage", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	// TODO: handlers

	// TODO: build grpc app

	// TODO: grpc app

	slog.Info("all is ok!")
}
