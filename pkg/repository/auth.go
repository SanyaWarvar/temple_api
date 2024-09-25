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
	query := fmt.Sprintf("INSERT INTO %s (id, username, email, password_hash) VALUES ($1, $2, $3, $4)", usersTable)
	_, err := r.db.Exec(query, id, user.Username, user.Email, user.Password)
	return err
}

func (r *AuthPostgres) GetUser(username, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username = $1 AND password_hash = $2", usersTable)
	err := r.db.Get(&user, query, username, password)
	return user, err
}

func (r *AuthPostgres) CheckEmailConfirm(email string) (bool, error) {
	var status bool
	query := fmt.Sprintf("SELECT confirmed_email FROM %s WHERE email = $1", usersTable)
	err := r.db.Get(&status, query, email)
	return status, err
}

func (r *AuthPostgres) ConfirmEmail(email string) error {

	query := fmt.Sprintf("UPDATE %s SET confirmed_email=true WHERE email = $1 ", usersTable)
	_, err := r.db.Exec(query, email)
	return err
}
