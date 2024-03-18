package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rmntim/movielab/internal/entity"
	"github.com/rmntim/movielab/internal/storage"
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

	query := fmt.Sprintf(
		`SELECT m.*, array_to_json(array_agg(a)) FROM movies m
				JOIN movie_actors ma ON ma.movie_id = m.id
				JOIN actors a ON a.id = ma.actor_id
				GROUP BY m.id
				ORDER BY $1 %s LIMIT $2 OFFSET $3`,
		orderDir)
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
		var movie entity.Movie
		err = rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating, &movie.Actors)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func (s *Storage) GetMovieById(id int) (*entity.Movie, error) {
	const op = "storage.postgres.GetMovieById"

	stmt, err := s.db.Prepare(
		`SELECT m.*, array_to_json(array_agg(a)) FROM movies m
				JOIN movie_actors ma ON m.id = ma.movie_id
				JOIN actors a ON ma.actor_id = a.id
				WHERE m.id = $1
				GROUP BY m.id`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var movie entity.Movie
	err = stmt.QueryRow(id).Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating, &movie.Actors)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrMovieNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &movie, nil
}
