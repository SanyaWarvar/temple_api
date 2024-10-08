package models

import (
	"regexp"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID
	Username string `json:"username" binding:"required" db:"username"`
	Email    string `json:"email" binding:"required" db:"email"`
	Password string `json:"password" binding:"required" db:"password_hash"`
}

const (
	UsernamePattern = "^[-a-zA-Z0-9_#$&*]+$"
	PasswordPattern = "^[-a-zA-Z0-9_#$&*]+$"
	UsernameMaxLen  = 32
	UsernameMinLen  = 4
	PasswordMaxLen  = 32
	PasswordMinLen  = 8
)

func (u *User) IsValid() bool {
	matched, err := regexp.Match(UsernamePattern, []byte(u.Username))
	usernameLen := len([]rune(u.Username))
	passwordLen := len([]rune(u.Password))
	if err != nil || !matched {
		return false
	}

	matched, err = regexp.Match(PasswordPattern, []byte(u.Password))
	if err != nil || !matched {
		return false
	}

	if (usernameLen <= UsernameMaxLen && usernameLen >= UsernameMinLen) &&
		(passwordLen <= PasswordMaxLen && passwordLen >= PasswordMinLen) {
		return true
	}
	return false

}

type UserInfo struct {
	UserId     uuid.UUID  `json:"-" db:"user_id"`
	FirstName  *string    `json:"first_name" db:"first_name"`
	SecondName *string    `json:"second_name" db:"second_name"`
	Status     *string    `json:"status" db:"status"`
	Birthday   *time.Time `json:"birthday" db:"birthday"`
	Gender     *string    `json:"gender" db:"gender"`
	Country    *string    `json:"country" db:"country"`
	City       *string    `json:"city" db:"city"`
}
