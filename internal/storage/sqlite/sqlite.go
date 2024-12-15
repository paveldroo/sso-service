package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"

	"github.com/paveldroo/sso-service/internal/config"
	"github.com/paveldroo/sso-service/internal/domain/models"
	"github.com/paveldroo/sso-service/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(cfg *config.Config) (*Storage, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("open database connection: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) AddUser(ctx context.Context, email string, pass_hash []byte) (int64, error) {
	stmt, err := s.db.Prepare("INSERT INTO users (email, pass_hash) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("prepare add user sql statement: %w", err)
	}

	res, err := stmt.ExecContext(ctx, email, pass_hash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("execute add user sql statement: %w", storage.ErrUserExists)
		}
		return 0, fmt.Errorf("execute add user sql statement: %w", err)
	}

	userID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get user id from sql result: %w", err)
	}

	return userID, nil
}

func (s *Storage) User(ctx context.Context, email string, pass_hash []byte) (models.User, error) {
	stmt, err := s.db.Prepare("SELECT id FROM users WHERE email=? AND password=?")
	if err != nil {
		return models.User{}, fmt.Errorf("prepare get user sql statement: %w", err)
	}

	row := stmt.QueryRowContext(ctx, email, pass_hash)

	var user models.User
	if err := row.Scan(&user.ID, &user.Email, &user.PassHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("execute get user sql statement: %w", storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("execute get user sql statement: %w", err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	stmt, err := s.db.PrepareContext(ctx, "SELECT is_admin FROM users WHERE id=?")
	if err != nil {
		return false, fmt.Errorf("prepare is_admin sql statement: %w", err)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool

	if err := row.Scan(&isAdmin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("execute is_admin sql statement: %w", storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("execute is_admin sql statement: %w", err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int64) (models.App, error) {
	stmt, err := s.db.Prepare("SELECT id FROM apps WHERE id=?")
	if err != nil {
		return models.App{}, fmt.Errorf("prepare get app sql statement: %w", err)
	}

	row := stmt.QueryRowContext(ctx, appID)

	var app models.App
	if err := row.Scan(&app.ID, &app.Name, &app.Secret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("execute get app sql statement: %w", storage.ErrUserNotFound)
		}
		return models.App{}, fmt.Errorf("execute get app sql statement: %w", err)
	}

	return app, nil
}
