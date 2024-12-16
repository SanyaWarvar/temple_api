package repository

import (
	"errors"
	"fmt"

	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TiktokPostgres struct {
	db *sqlx.DB
}

func NewTiktokPostgres(db *sqlx.DB) *TiktokPostgres {
	return &TiktokPostgres{db: db}
}

func (r *TiktokPostgres) CreateTiktok(item models.Tiktok) error {
	query := fmt.Sprintf(
		`
		INSERT INTO %s VALUES ($1, $2, $3, $4, $5)
		`,
		tiktoksTable,
	)
	_, err := r.db.Exec(query, item.Id, item.AuthorId, item.CreatedAt, item.Title, item.Body)
	return err
}

func (r *TiktokPostgres) GetTiktokById(tiktokId uuid.UUID) (models.Tiktok, error) {
	var item models.Tiktok
	query := fmt.Sprintf(
		`
		SELECT id, (select username from users where id = t.author_id), created_at, title, id as body FROM %s t WHERE id = $1
		`,
		tiktoksTable,
	)
	err := r.db.Get(&item, query, tiktokId)
	return item, err
}

func (r *TiktokPostgres) DeleteTiktokById(tiktokId, userId uuid.UUID) error {

	query := fmt.Sprintf(
		`
		DELETE FROM %s WHERE id = $1 AND author_id = $2
		`,
		tiktoksTable,
	)
	res, err := r.db.Exec(query, tiktokId, userId)
	count, _ := res.RowsAffected()
	if count == 0 {
		return errors.New("bad video id")
	}
	return err
}
