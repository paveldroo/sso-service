package sqlite

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func New() error {
	if err := migrateDB(); err != nil {
		return err
	}

	return nil
}

func migrateDB() error {
	m, err := migrate.New(
		"file://../../internal/db/sqlite/migrations",
		"sqlite3://../../internal/db/sqlite/db.sqlite",
	)

	if err != nil {
		return fmt.Errorf("init migrations: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
