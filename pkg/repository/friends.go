package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FriendPostgres struct {
	db *sqlx.DB
}

func NewFriendPostgres(db *sqlx.DB) *FriendPostgres {
	return &FriendPostgres{db: db}
}

type FriendOutput struct {
	Fullname      string `json:"fullname" db:"fullname"`
	Username      string `json:"username" db:"username"`
	ProfilePicUrl string `json:"profile_pciture_url" db:"profile_pciture"`
}

type FriendListOutput struct {
	Friends []FriendOutput `json:"friends" db:"friends"`
}

func (r *FriendPostgres) GetAllFriend(userId uuid.UUID, page int) (FriendListOutput, error) {
	offset := (page - 1) * 50
	query := fmt.Sprintf(
		`
		SELECT 
			ui.first_name || ' ' || ui.second_name AS full_name,
			u.username,
			ui.profile_picture
		FROM 
			%s fi
		JOIN 
			%s u ON u.id = fi.from_user_id OR u.id = fi.to_user_id
		JOIN 
			%s ui ON ui.user_id = u.id
		WHERE 
			(fi.from_user_id = $1 OR fi.to_user_id = $1) 
			AND fi.confirmed = 't'
			AND u.id != $1
		OFFSET $2
		LIMIT 50
		`,
		friendsTable,
		usersTable,
		usersInfoTable,
	)

	rows, err := r.db.Query(query, userId, offset)
	if err != nil {
		return FriendListOutput{}, err
	}
	defer rows.Close()

	var friends []FriendOutput
	for rows.Next() {
		var friend FriendOutput
		if err := rows.Scan(&friend.Fullname, &friend.Username, &friend.ProfilePicUrl); err != nil {
			return FriendListOutput{}, err
		}
		friends = append(friends, friend)
	}

	if err := rows.Err(); err != nil {
		return FriendListOutput{}, err
	}

	return FriendListOutput{Friends: friends}, nil
}

func (r *FriendPostgres) InviteFriend(fromId uuid.UUID, toUsername string) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (from_user_id, to_user_id) 
		VALUES ($1, (SELECT id FROM %s WHERE username = $2))`,
		friendsTable,
		usersTable,
	)
	_, err := r.db.Exec(query, fromId, toUsername)
	return err
}

func (r *FriendPostgres) DeleteByU(invitedId uuid.UUID, ownerUsername string) error {
	query := fmt.Sprintf(
		`UPDATE %s
		SET confirmed = 'f'
		WHERE from_user_id = (SELECT id FROM %s WHERE username = $1) AND to_user_id = $2`,
		friendsTable,
		usersTable,
	) //если ты удаляешь из друзей того, кто тебе отправлял запрос

	res, err := r.db.Exec(query, ownerUsername, invitedId)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 1 {
		return nil
	}

	query = fmt.Sprintf(
		`DELETE FROM %s WHERE from_user_id = $1 AND to_user_id = (SELECT id FROM %s WHERE username = $2)`,
		friendsTable,
		usersTable,
	) //если ты удаляешь из друзей того, на кого ты подписан

	res, err = r.db.Exec(query, invitedId, ownerUsername)
	if err != nil {
		fmt.Println(err)
		return err
	}
	n, err = res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("%s not in your friends", ownerUsername)
	}
	return nil
}

func (r *FriendPostgres) ConfirmFriend(invitedId uuid.UUID, ownerUsername string) error {
	query := fmt.Sprintf(
		`UPDATE %s
		SET confirmed = 't'
		WHERE from_user_id = (SELECT id FROM %s WHERE username = $1) AND to_user_id = $2`,
		friendsTable,
		usersTable,
	)
	res, err := r.db.Exec(query, ownerUsername, invitedId)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf("%s not invite you", ownerUsername)
	}

	return err
}
