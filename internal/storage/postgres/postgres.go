package postgres

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rmntim/movielab/internal/entity"
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

func (s *Storage) GetMovies(limit, offset int, orderBy string, asc bool) ([]entity.Movie, error) {
	const op = "storage.postgres.GetMovies"

	orderDir := "DESC"
	if asc {
		orderDir = "ASC"
	}
	// HACK: pretty funny way of querying json data, but it probably more performant
	// than querying actors for each movie or getting every actor/movie record in memory
	query := fmt.Sprintf("SELECT row_to_json(row) FROM (SELECT * FROM movies_actors ORDER BY $1 %s LIMIT $2 OFFSET $3) row", orderDir)
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.Query(orderBy, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var movies []entity.Movie
	for rows.Next() {
		var movieJson string
		err = rows.Scan(&movieJson)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		var movie entity.Movie
		if err := json.Unmarshal([]byte(movieJson), &movie); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		movies = append(movies, movie)
	}

	return movies, nil
}
