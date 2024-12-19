package repository

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MessagesPostgres struct {
	db *sqlx.DB
}

func NewMessagesPostgres(db *sqlx.DB) *MessagesPostgres {
	return &MessagesPostgres{db: db}
}

func (r *MessagesPostgres) CreateChat(inviteUsername string, owner uuid.UUID) (uuid.UUID, error) {
	chatId := uuid.New()
	tx, err := r.db.Begin()
	if err != nil {
		return chatId, err
	}
	query := fmt.Sprintf(`
		INSERT INTO %s VALUES ($1)
	`, chatsTable)
	_, err = tx.Exec(query, chatId)
	if err != nil {
		return chatId, err
	}

	query = fmt.Sprintf(`
		INSERT INTO %s SELECT $1, id FROM %s WHERE username = $2 OR id = $3
	`, chatMembersTable, usersTable)

	_, err = tx.Exec(query, chatId, inviteUsername, owner)

	if err != nil {
		return chatId, err
	}
	err = tx.Commit()

	return chatId, err
}

type MessageOutput struct {
	Id               uuid.UUID `json:"id" db:"id"`
	Body             string    `json:"body" db:"body"`
	AuthorFirstName  string    `json:"author_first_name" db:"author_first_name"`
	AuthorSecondName string    `json:"author_second_name" db:"author_second_name"`
	AuthorProfilePic string    `json:"author_profile_picture" db:"author_profile_picture"`
	ChatId           uuid.UUID `json:"chat_id" db:"chat_id"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	Readed           bool      `json:"readed" db:"readed"`
	Edited           bool      `json:"edited" db:"edited"`
}

type WithUserStruct struct {
	FirstName  string `json:"first_name" db:"first_name"`
	SecondName string `json:"second_name" db:"second_name"`
	ProfilePic string `json:"profile_picture" db:"profile_picture"`
	Username   string `json:"username" db:"username"`
}

type AllChatsOutput struct {
	ChatId   uuid.UUID       `json:"chat_id" db:"chat_id"`
	WithUser WithUserStruct  `json:"with_user" db:"with_user"`
	Messages []MessageOutput `json:"messages" db:"messages"`
}

func (r *MessagesPostgres) GetAllChats(userId uuid.UUID, page int) ([]AllChatsOutput, error) {
	offset := (page - 1) * 25
	query := `WITH user_chats AS (
		SELECT c.id AS chat_id
		FROM chats c
		JOIN chat_members cm ON c.id = cm.chat_id
		WHERE cm.user_id = $1
		LIMIT 25 OFFSET $2
	),
	messages_with_row_number AS (
		SELECT 
			m.*,
			ROW_NUMBER() OVER (PARTITION BY m.chat_id ORDER BY m.created_at DESC) AS rn
		FROM messages m
	)
	SELECT 
		uc.chat_id,
		(SELECT username FROM users WHERE id = mw.author_id) AS username,
		(SELECT first_name FROM users_info WHERE user_id = cm.user_id) AS with_user_first_name,
		(SELECT second_name FROM users_info WHERE user_id = cm.user_id) AS with_user_second_name,
		(SELECT profile_picture FROM users_info WHERE user_id = cm.user_id) AS with_user_profile_picture,
		mw.id,
		mw.body,
		(SELECT first_name FROM users_info WHERE user_id = mw.author_id) AS message_author_first_name,
		(SELECT second_name FROM users_info WHERE user_id = mw.author_id) AS message_author_second_name,
		(SELECT profile_picture FROM users_info WHERE user_id = mw.author_id) AS message_author_profile_picture,
		mw.created_at AS message_created_at,
		mw.readed AS message_readed,
		mw.edited AS message_edited,
		mw.chat_id
	FROM user_chats uc
	RIGHT JOIN chat_members cm ON cm.chat_id = uc.chat_id
	RIGHT JOIN messages_with_row_number mw ON mw.chat_id = uc.chat_id AND mw.rn <= 50
	WHERE cm.user_id != $1
	ORDER BY uc.chat_id, mw.created_at`

	rows, err := r.db.Queryx(query, userId, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chatsMap := make(map[uuid.UUID]*AllChatsOutput)

	for rows.Next() {
		var chatId uuid.UUID
		var withUser WithUserStruct
		var message MessageOutput

		err := rows.Scan(
			&chatId,
			&withUser.Username,
			&withUser.FirstName,
			&withUser.SecondName,
			&withUser.ProfilePic,
			&message.Id,
			&message.Body,
			&message.AuthorFirstName,
			&message.AuthorSecondName,
			&message.AuthorProfilePic,
			&message.CreatedAt,
			&message.Readed,
			&message.Edited,
			&message.ChatId,
		)
		if err != nil {
			return nil, err
		}

		if message.Body == "" {
			continue
		}

		if chat, exists := chatsMap[chatId]; exists {
			chat.Messages = append(chat.Messages, message)
		} else {
			chatsMap[chatId] = &AllChatsOutput{
				ChatId:   chatId,
				WithUser: withUser,
				Messages: []MessageOutput{message},
			}
		}
	}

	var chats []AllChatsOutput
	for _, chat := range chatsMap {
		chats = append(chats, *chat)
	}

	return chats, nil
}

func (r *MessagesPostgres) GetChat(chatId, userId uuid.UUID, page int) (AllChatsOutput, error) {
	// Определяем смещение для пагинации
	offset := (page - 1) * 100

	query := `
	SELECT 
		uc.id,
		(SELECT username FROM users WHERE id = mw.author_id) AS username,
		(SELECT first_name FROM users_info WHERE user_id = cm.user_id) AS with_user_first_name,
		(SELECT second_name FROM users_info WHERE user_id = cm.user_id) AS with_user_second_name,
		(SELECT profile_picture FROM users_info WHERE user_id = cm.user_id) AS with_user_profile_picture,
		mw.id,
		mw.body,
		(SELECT first_name FROM users_info WHERE user_id = mw.author_id) AS message_author_first_name,
		(SELECT second_name FROM users_info WHERE user_id = mw.author_id) AS message_author_second_name,
		(SELECT profile_picture FROM users_info WHERE user_id = mw.author_id) AS message_author_profile_picture,
		mw.created_at AS message_created_at,
		mw.readed AS message_readed,
		mw.edited AS message_edited,
		mw.chat_id
	FROM chats uc
	RIGHT JOIN chat_members cm ON cm.chat_id = uc.id
	RIGHT JOIN messages mw ON mw.chat_id = uc.id
	WHERE cm.user_id != $1 AND uc.id = $2
	ORDER BY uc.id, mw.created_at
	limit 100 offset $3`

	rows, err := r.db.Queryx(query, userId, chatId, offset)
	if err != nil {
		return AllChatsOutput{}, err
	}
	defer rows.Close()

	var chat AllChatsOutput
	var withUser WithUserStruct

	chat.ChatId = chatId

	// Сбор участников и сообщений
	for rows.Next() {
		var message MessageOutput

		err := rows.Scan(
			&chatId,
			&withUser.Username,
			&withUser.FirstName,
			&withUser.SecondName,
			&withUser.ProfilePic,
			&message.Id,
			&message.Body,
			&message.AuthorFirstName,
			&message.AuthorSecondName,
			&message.AuthorProfilePic,
			&message.CreatedAt,
			&message.Readed,
			&message.Edited,
			&message.ChatId,
		)
		if err != nil {
			return AllChatsOutput{}, err
		}

		// Добавляем сообщение в чат
		chat.Messages = append(chat.Messages, message)
	}

	chat.WithUser = withUser

	return chat, nil

}

func (r *MessagesPostgres) CreateMessage(data models.Message) error {
	var err error
	if data.ReplyTo.String() == "00000000-0000-0000-0000-000000000000" {
		query := fmt.Sprintf(
			`
			INSERT INTO %s (id, body, author_id, chat_id) VALUES ($1,$2, $3, $4)
			`, messagesTable,
		)
		_, err = r.db.Exec(query, data.Id, data.Body, data.AuthorId, data.ChatId)
	} else {
		query := fmt.Sprintf(
			`
			INSERT INTO %s (id, body, author_id, chat_id, reply_to) VALUES ($1,$2, $3, $4, $5)
			`, messagesTable,
		)
		_, err = r.db.Exec(query, data.Id, data.Body, data.AuthorId, data.ChatId, data.ReplyTo)
	}

	return err
}

func (r *MessagesPostgres) ReadMessage(messageId, userId uuid.UUID) error {
	query := fmt.Sprintf( // добавить проверку что прочитал пользователь который внутри чата?
		`
		UPDATE %s SET readed = 't' WHERE id = $1 AND author_id != $2 
		`, messagesTable,
	)
	rows, err := r.db.Exec(query, messageId, userId)
	count, _ := rows.RowsAffected()
	if count == 0 {
		err = fmt.Errorf("message not found or you cant read your own message")
	}
	return err
}

func (r *MessagesPostgres) EditMessage(userId uuid.UUID, message models.Message) error {
	fields := make([]string, 0)
	values := make([]interface{}, 0)

	rv := reflect.ValueOf(message)
	trueCounter := 1
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			fieldName := rv.Type().Field(i).Tag.Get("db")
			fields = append(fields, fmt.Sprintf("%s=$%d", fieldName, trueCounter))
			values = append(values, field.Elem().Interface())
			trueCounter += 1
		}
	}

	if len(fields) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE messages SET %s WHERE id = $%d AND author_id = $%d",
		strings.Join(fields, ", "), trueCounter, trueCounter+1)

	values = append(values, message.Id, userId)

	_, err := r.db.Exec(query, values...)
	return err
}

func (r *MessagesPostgres) DeleteMessage(messageId, userId uuid.UUID) error {
	query := fmt.Sprintf(
		`
		DELETE FROM %s WHERE id = $1 AND author_id = $2 
		`, messagesTable,
	)
	rows, err := r.db.Exec(query, messageId, userId)
	count, _ := rows.RowsAffected()
	if count == 0 {
		err = fmt.Errorf("message not found or you cant delete this message")
	}
	return err
}

func (r *MessagesPostgres) GetMembersFromChatByID(chatId uuid.UUID) ([]models.User, error) {
	var output []models.User
	query := fmt.Sprintf(
		`
		SELECT u.username
		FROM %s cm 
		LEFT JOIN %s c on c.id = cm.chat_id
		LEFT JOIN %s u on u.id = cm.user_id
		WHERE cm.chat_id = $1
		`, chatMembersTable, chatsTable, usersTable,
	)
	err := r.db.Select(&output, query, chatId)
	return output, err
}
