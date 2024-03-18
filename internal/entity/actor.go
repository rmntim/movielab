package entity

import "time"

type Actor struct {
	ID int `json:"id"`
	NewActor
	MovieIDs []int32 `json:"movie_ids"`
}

type NewActor struct {
	Name      string    `json:"name"`
	Sex       string    `json:"sex"`
	BirthDate time.Time `json:"birthdate"`
}
