package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/paveldroo/sso-service/internal/config"
	server "github.com/paveldroo/sso-service/internal/grpc"
	slogpretty "github.com/paveldroo/sso-service/internal/lib/logger/handlers"
	"github.com/paveldroo/sso-service/internal/lib/logger/sl"
	"github.com/paveldroo/sso-service/internal/storage/sqlite"
)

const envLocal = "local"

func main() {
	cfg := config.MustLoad()

	setupLogger(cfg.Env)

	storage, err := sqlite.New(cfg)
	if err != nil {
		slog.Error("create sqlite storage", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	s := server.New(storage)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	slog.Info("starting server...")
	if err = s.Serve(lis); err != nil {
		slog.Error("running server", sl.Err(err))
		os.Exit(1)
	}

	slog.Info("all is ok!")
}

func setupLogger(env string) {
	if env == envLocal {
		setupPrettySlog()
		return
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(jsonHandler))
}

func setupPrettySlog() {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)
	slog.SetDefault(slog.New(handler))
}
