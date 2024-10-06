package repository

import (
	"fmt"

	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) CreateUser(user models.User) error {
	id := uuid.NewString()
	query := fmt.Sprintf(`INSERT INTO %s (id, username, email, password_hash) VALUES ($1, $2, $3, $4)`, usersTable)
	_, err := r.db.Exec(query, id, user.Username, user.Email, user.Password)
	return err
}

func (r *UserPostgres) GetUserByUP(username, hashedPassword string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf(`SELECT id FROM %s WHERE username = $1 AND password_hash = $2`, usersTable)
	err := r.db.Get(&user, query, username, hashedPassword)
	return user, err
}

func (r *UserPostgres) GetUserByU(username string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf(`SELECT id, username, email, password_hash FROM %s WHERE username = $1`, usersTable)
	err := r.db.Get(&user, query, username)
	return user, err
}

func (r *UserPostgres) GetUserById(userId uuid.UUID) (models.User, error) {
	var user models.User
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, usersTable)
	err := r.db.Get(&user, query, userId)
	return user, err
}

func (m *UserPostgres) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func (r *UserPostgres) ComparePassword(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
