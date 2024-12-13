package main

import (
	"log/slog"

	"github.com/paveldroo/sso-service/internal/db/sqlite"
	"github.com/paveldroo/sso-service/internal/logger/sl"
)

func main() {
	err := sqlite.New()
	if err != nil {
		slog.Error("run migrations", sl.Err(err))

	}
	// TODO: config

	// TODO: logger

	// TODO: storage

	// TODO: handlers

	// TODO: build grpc app

	// TODO: grpc app
}
