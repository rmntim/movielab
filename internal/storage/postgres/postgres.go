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

func (s *Storage) DeleteMovie(id int) error {
	const op = "storage.postgres.DeleteMovie"

	stmt, err := s.db.Prepare("DELETE FROM movies WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateMovie(id int, movie *entity.Movie) error {
	const op = "storage.postgres.UpdateMovie"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("UPDATE movies SET title = $1, description = $2, release_date = $3, rating = $4 WHERE id = $5")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(movie.Title, movie.Description, movie.ReleaseDate, movie.Rating, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = tx.Prepare("DELETE FROM movie_actors WHERE movie_id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = tx.Prepare("INSERT INTO movie_actors (movie_id, actor_id) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, actorID := range movie.ActorIDs {
		_, err = stmt.Exec(id, actorID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SearchMovies(title, actorName string, limit, offset int) ([]entity.Movie, error) {
	const op = "storage.postgres.SearchMovies"

	stmt, err := s.db.Prepare(
		`SELECT m.*, array_remove(array_agg(a.id), NULL) FROM movies m
				LEFT JOIN movie_actors ma ON ma.movie_id = m.id
				LEFT JOIN actors a ON a.id = ma.actor_id
				WHERE m.title ILIKE $1 AND a.name ILIKE $2
				GROUP BY m.id
				LIMIT $3 OFFSET $4`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.Query(fmt.Sprintf("%%%s%%", title), fmt.Sprintf("%%%s%%", actorName), limit, offset)
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

func (s *Storage) GetActors(limit, offset int) ([]entity.Actor, error) {
	const op = "storage.postgres.GetActors"

	stmt, err := s.db.Prepare(
		`SELECT a.*, array_remove(array_agg(m.id), NULL) FROM actors a
				LEFT JOIN movie_actors ma ON ma.actor_id = a.id
				LEFT JOIN movies m ON m.id = ma.movie_id
				GROUP BY a.id
				LIMIT $1 OFFSET $2`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.Query(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var actors []entity.Actor
	for rows.Next() {
		var actor entity.Actor
		err = rows.Scan(&actor.ID, &actor.Name, &actor.Sex, &actor.BirthDate, (*pq.Int32Array)(&actor.MovieIDs))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		actors = append(actors, actor)
	}

	return actors, nil
}

func (s *Storage) GetActorById(id int) (*entity.Actor, error) {
	const op = "storage.postgres.GetActorByID"

	stmt, err := s.db.Prepare(`
		SELECT a.*, array_remove(array_agg(m.id), NULL) FROM actors a
		LEFT JOIN movie_actors ma ON ma.actor_id = a.id
		LEFT JOIN movies m ON m.id = ma.movie_id
		WHERE a.id = $1
		GROUP BY a.id`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var actor entity.Actor
	err = stmt.QueryRow(id).Scan(&actor.ID, &actor.Name, &actor.Sex, &actor.BirthDate, (*pq.Int32Array)(&actor.MovieIDs))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrActorNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &actor, nil
}

func (s *Storage) CreateActor(actor *entity.NewActor) (int, error) {
	const op = "storage.postgres.CreateActor"

	stmt, err := s.db.Prepare("INSERT INTO actors (name, sex, birth_date) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int
	err = stmt.QueryRow(actor.Name, actor.Sex, actor.BirthDate).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) DeleteActor(id int) error {
	const op = "storage.postgres.DeleteActor"

	stmt, err := s.db.Prepare("DELETE FROM actors WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateActor(id int, actor *entity.Actor) error {
	const op = "storage.postgres.UpdateActor"

	stmt, err := s.db.Prepare("UPDATE actors SET name = $1, sex = $2, birth_date = $3 WHERE id = $4")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(actor.Name, actor.Sex, actor.BirthDate, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
