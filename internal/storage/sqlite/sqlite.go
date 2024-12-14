package sqlite

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/paveldroo/sso-service/internal/config"
)

func New(cfg *config.Config) error {
	if err := migrateDB(cfg); err != nil {
		return err
	}

	return nil
}

func migrateDB(cfg *config.Config) error {
	m, err := migrate.New(
		"file://"+cfg.DB.MigrationsPath,
		"sqlite3://"+cfg.DB.Path,
	)

	if err != nil {
		return fmt.Errorf("init migrations: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
