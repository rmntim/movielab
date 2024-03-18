package entity

import "time"

type Actor struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Sex       string    `json:"sex"`
	BirthDate time.Time `json:"birthdate" db:"birth_date"`
}
