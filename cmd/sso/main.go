package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/paveldroo/sso-service/internal/app"
	"github.com/paveldroo/sso-service/internal/config"
	slogpretty "github.com/paveldroo/sso-service/internal/lib/logger/handlers"
)

const envLocal = "local"

func main() {
	cfg := config.MustLoad()

	setupLogger(cfg.Env)

	app := app.New(slog.Default(), cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTLL)

	go func() {
		app.GRPCServer.MustRun()
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.GRPCServer.Stop()
	slog.Info("server stopped")
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
