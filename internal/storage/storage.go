package storage

import (
	"errors"
)

var (
	ErrMovieNotFound = errors.New("movie not found")

	ErrActorNotFound = errors.New("actor not found")
)
