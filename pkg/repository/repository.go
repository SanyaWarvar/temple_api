package repository

import (
	"time"

	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Authorizer interface {
	CreateUser(user models.User) error
	GetUser(username, password_hash string) (models.User, error)
	CheckEmailConfirm(email string) (bool, error)
	ConfirmEmail(email string) error
}

type Cacher interface {
	GetConfirmCode(email string) (string, time.Duration, error)
	SaveConfirmCode(email, code string) error
}

type Repository struct {
	Authorizer
	Cacher
}

func NewRepository(db *sqlx.DB, cacheDb *redis.Client, codeExp time.Duration) *Repository {
	return &Repository{
		Authorizer: NewAuthPostgres(db),
		Cacher:     NewCache(cacheDb, codeExp),
	}
}
