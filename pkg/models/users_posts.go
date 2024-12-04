package models

import (
	"time"

	"github.com/google/uuid"
)

type UserPost struct {
	Id         uuid.UUID `json:"id" db:"id"`
	AuthorId   uuid.UUID `json:"author_id" db:"author_id"`
	Body       string    `json:"body" db:"body"`
	LastUpdate time.Time `json:"last_update" db:"last_update"`
	Edited     bool      `json:"edited" db:"edited"`
}
