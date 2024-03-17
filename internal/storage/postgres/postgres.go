package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sqlx.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
func (s *Storage) GetUserRole(username string, password string) (string, error) {
	const op = "storage.postgres.GetUserRole"

	stmt, err := s.db.Prepare("SELECT role FROM users WHERE username = $1 AND password = $2")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var role string
	err = stmt.QueryRow(username, password).Scan(&role)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}
