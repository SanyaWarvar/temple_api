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

type TikTokOutput struct {
	models.Tiktok
	AuthorFirstName      string `json:"author_first_name" db:"author_first_name"`
	AuthorSecondName     string `json:"author_second_name" db:"author_second_name"`
	AuthorProfilePicture string `json:"author_profile_picture" db:"author_profile_picture"`
}

func (r *TiktokPostgres) Feed(userId uuid.UUID, page int) ([]TikTokOutput, error) {
	offset := (page - 1) * 50
	var output []TikTokOutput
	query := fmt.Sprint(
		`
		select 
		t.*,
		u.username as username,
		ui.first_name as author_first_name,
		ui.second_name as author_second_name,
		ui.profile_picture as author_profile_picture

		from tiktoks t 
		inner join users u on u.id = t.author_id
		inner join users_info ui on ui.user_id = t.author_id 
		order by created_at desc
		limit 50 offset $1
		`,
	)
	err := r.db.Select(&output, query, offset)

	return output, err
}
