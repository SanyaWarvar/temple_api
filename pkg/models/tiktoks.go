package models

import (
	"time"

	"github.com/google/uuid"
)

type Tiktok struct {
	Id        uuid.UUID `json:"id" db:"id"`
	AuthorId  uuid.UUID `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Title     string    `json:"title" db:"title"`
	Body      string    `json:"body" db:"body"`
}
