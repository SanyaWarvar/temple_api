package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Id        uuid.UUID  `json:"id" db:"id"`
	Body      *string    `json:"body" db:"body"`
	AuthorId  uuid.UUID  `json:"author_id" db:"author_id"`
	ChatId    uuid.UUID  `json:"chat_id" db:"chat_id"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	Readed    *bool      `json:"readed" db:"readed"`
	Edited    *bool      `json:"edited" db:"edited"`
	ReplyTo   uuid.UUID  `json:"reply_to" db:"reply_to"`
}

func (m *Message) IsValid() bool {
	if m.Body == nil ||
		m.AuthorId.String() == "00000000-0000-0000-0000-000000000000" ||
		m.ChatId.String() == "00000000-0000-0000-0000-000000000000" {
		return false
	}
	return true

}
