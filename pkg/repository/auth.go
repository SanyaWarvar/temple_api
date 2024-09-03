package repository

import (
	"fmt"

	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user models.User) error {
	id := uuid.NewString()
	query := fmt.Sprintf("INSERT INTO %s (id, username, password_hash) VALUES ($1, $2, $3)", usersTable)
	_, err := r.db.Exec(query, id, user.Username, user.Password)
	return err
}

func (r *AuthPostgres) GetUser(username, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username = $1 AND password_hash = $2", usersTable)
	err := r.db.Get(&user, query, username, password)
	return user, err
}
