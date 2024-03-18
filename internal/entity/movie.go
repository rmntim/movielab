package entity

import "time"

// Movie represents a movie record
type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      int       `json:"rating"`
}
