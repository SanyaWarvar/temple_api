package models

import "github.com/google/uuid"

type Chat struct {
	Id       uuid.UUID `json:"id" db:"id"`
	Members  []string  `json:"members"`
	Messages []Message `json:"messages"`
}
