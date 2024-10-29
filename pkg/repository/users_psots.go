package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UsersPostsPostgres struct {
	db *sqlx.DB
}

func NewUsersPostsPostgres(db *sqlx.DB) *UsersPostsPostgres {
	return &UsersPostsPostgres{db: db}
}

func (r *UsersPostsPostgres) CreatePost(post models.UserPost) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (id, author_id, body, last_update) VALUES($1, $2, $3, $4)`,
		usersPostsTable,
	)
	_, err := r.db.Exec(query, post.Id, post.AuthorId, post.Body, post.LastUpdate)
	return err
}
func (r *UsersPostsPostgres) UpdatePost(newPost models.UserPost) error {
	query := fmt.Sprintf(`UPDATE %s SET body = $1, edited = 't', last_update = $2 WHERE id = $3`, usersPostsTable)
	res, err := r.db.Exec(query, newPost.Body, newPost.LastUpdate, newPost.Id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("post not found")
	}
	return err
}
func (r *UsersPostsPostgres) GetPostById(postId, userId uuid.UUID) (UserPostOutput, error) {
	var post UserPostOutput
	query := fmt.Sprintf(`
		SELECT 
		p.id, p.author_id, p.body, p.last_update, p.edited, 
		count(p.id) as likes_count,
		(select count (*) from %s upl2 where upl2.post_id = p.id and upl2.user_id = $2) as liked_by_me 
		FROM %s p 
		INNER JOIN %s upl on upl.post_id = p.id 
		WHERE p.id = $1
		GROUP BY p.id, p.author_id, p.body, p.last_update, p.edited 
		
	`, usersPostsLikesTable, usersPostsTable, usersPostsLikesTable)
	err := r.db.Get(&post, query, postId, userId)
	return post, err
}
func (r *UsersPostsPostgres) DeletePostById(postId, userId uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1 AND author_id = $2`, usersPostsTable)
	res, err := r.db.Exec(query, postId, userId)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("post not found")
	}
	return err
}

type UserPostOutput struct {
	Id         uuid.UUID `json:"id" db:"id"`
	AuthorId   uuid.UUID `json:"author_id" db:"author_id"`
	Body       string    `json:"body" db:"body"`
	LastUpdate time.Time `json:"last_update" db:"last_update"`
	Edited     bool      `json:"edited" db:"edited"`
	LikesCount int       `json:"likes_count" db:"likes_count"`
	LikedByMe  bool      `json:"liked_by_me" db:"liked_by_me"`
}

func (r *UsersPostsPostgres) GetPostsByU(username string, page int, userId uuid.UUID) ([]UserPostOutput, error) {
	data := make([]UserPostOutput, 0, 50)
	query := fmt.Sprintf(`
		with posts as (
		select * from %s up where author_id = (select id from %s where username = $1)
		)
		select p.id, p.author_id, p.body, p.last_update, p.edited, 
		count(p.id) as likes_count, 
		(select count (*) from %s upl2 where upl2.post_id = p.id and upl2.user_id = $2) as liked_by_me 
		from %s upl 
		inner join posts p on p.id = upl.post_id
		group by p.id, p.author_id, p.body, p.last_update, p.edited
		limit 50
		offset $3
		`,
		usersPostsTable,
		usersTable,
		usersPostsLikesTable,
		usersPostsLikesTable,
	)
	err := r.db.Select(&data, query, username, userId, page*50)
	return data, err
}

func (r *UsersPostsPostgres) LikePostById(postId, userId uuid.UUID) error {
	query := fmt.Sprintf(
		`
		INSERT INTO %s VALUES ($1, $2)
		`,
		usersPostsLikesTable,
	)
	_, err := r.db.Exec(query, postId, userId)

	if err != nil && strings.Contains(err.Error(), "pkey") {
		query = fmt.Sprintf(
			`
			DELETE FROM %s WHERE post_id = $1 AND user_id = $2
			`,
			usersPostsLikesTable,
		)
		_, err = r.db.Exec(query, postId, userId)
	}

	return err
}
