package entity

import (
	"encoding/json"
	"time"
)

// Movie represents a movie record
type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      int       `json:"rating"`
	Actors      ActorList `json:"actors"`
}

type ActorList []Actor

func (l *ActorList) Scan(src any) error {
	return json.Unmarshal(src.([]byte), l)
}
