package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/paveldroo/sso-service/internal/config"
	server "github.com/paveldroo/sso-service/internal/grpc"
	"github.com/paveldroo/sso-service/internal/logger/sl"
	"github.com/paveldroo/sso-service/internal/storage/sqlite"
)

func main() {
	cfg, err := config.ParseCfg("config.yml")
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

	s := server.New()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	slog.Info("starting server...")
	if err = s.Serve(lis); err != nil {
		slog.Error("running server: %w", err)
		os.Exit(1)
	}

	slog.Info("all is ok!")
}
