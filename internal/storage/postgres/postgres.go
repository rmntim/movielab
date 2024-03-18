package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
		`SELECT m.*, array_remove(array_agg(a.id), NULL) FROM movies m
				LEFT JOIN movie_actors ma ON ma.movie_id = m.id
				LEFT JOIN actors a ON a.id = ma.actor_id
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
		err = rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating, (*pq.Int32Array)(&movie.ActorIDs))
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
		`SELECT m.*, array_remove(array_agg(a.id), NULL) FROM movies m
				LEFT JOIN movie_actors ma ON m.id = ma.movie_id
				LEFT JOIN actors a ON ma.actor_id = a.id
				WHERE m.id = $1
				GROUP BY m.id`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var movie entity.Movie
	err = stmt.QueryRow(id).Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating, (*pq.Int32Array)(&movie.ActorIDs))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrMovieNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &movie, nil
}

func (s *Storage) CreateMovie(movie *entity.NewMovie) (int, error) {
	const op = "storage.postgres.CreateMovie"

	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO movies (title, description, release_date, rating) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int
	err = stmt.QueryRow(movie.Title, movie.Description, movie.ReleaseDate, movie.Rating).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = tx.Prepare("INSERT INTO movie_actors (movie_id, actor_id) VALUES ($1, $2)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	for _, actorID := range movie.ActorIDs {
		_, err = stmt.Exec(id, actorID)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
