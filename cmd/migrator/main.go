package main

import (
	"errors"
	"flag"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/paveldroo/sso-service/internal/lib/logger/sl"
)

func main() {
	storagePath := flag.String("storage-path", "", "path to storage file")
	migrationsPath := flag.String("migrations-path", "", "path to migrations folder")
	flag.Parse()

	if *storagePath == "" || *migrationsPath == "" {
		slog.Error("not all arguments are set up")
		os.Exit(1)
	}

	m, err := migrate.New(
		"file://"+*migrationsPath,
		"sqlite3://"+*storagePath,
	)
	if err != nil {
		slog.Error("failed to init migrations", sl.Err(err))
		os.Exit(1)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("no migrations to apply")
			return
		}
		slog.Error("failed to run migrations", sl.Err(err))
		os.Exit(1)
	}

	slog.Info("all migrations successfully applied")
}
