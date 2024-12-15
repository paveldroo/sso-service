package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"

	"github.com/paveldroo/sso-service/internal/config"
	"github.com/paveldroo/sso-service/internal/user"
)

type Sqlite struct {
	db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.DB.Path)
	if err != nil {
		return nil, fmt.Errorf("open database connection: %w", err)
	}

	return &Sqlite{db: db}, nil
}

func (s *Sqlite) AddUser(email, password string, isAdmin bool) error {
	q := fmt.Sprintf("INSERT INTO users (email, password, isAdmin) VALUES ('%s', '%s', '%v');", email, password, isAdmin)
	_, err := s.db.Exec(q)
	if err != nil {
		return fmt.Errorf("insert user in db: %w", err)
	}

	return nil
}

func (s *Sqlite) User(email, password string) (int64, error) {
	q := fmt.Sprintf("SELECT id FROM users WHERE email='%s' AND password='%s';", email, password)
	res, err := s.db.Exec(q)
	if err != nil {
		return 0, fmt.Errorf("select user from db: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id: %w", err)
	}

	return id, nil
}

func (s *Sqlite) IsAdmin(email string) (bool, error) {
	q := fmt.Sprintf("SELECT * FROM users WHERE email='%s';", email)
	rows, err := s.db.Query(q)
	if err != nil {
		return false, fmt.Errorf("select user from db: %w", err)
	}
	defer rows.Close()

	u := user.User{}
	for rows.Next() {
		if err = rows.Scan(&u.ID, &u.Email, &u.Password, &u.IsAdmin); err != nil {
			return false, fmt.Errorf("unmarshal db row to user struct: %w", err)
		}
	}

	return u.IsAdmin, nil
}
