package repository

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (r *MessagesPostgres) GetAllChats(userId uuid.UUID, page int) ([]models.Chat, error) {
	// Определяем смещение для пагинации
	offset := (page - 1) * 25

	query := `WITH user_chats AS (
                SELECT c.id AS chat_id
                FROM chats c
                JOIN chat_members cm ON c.id = cm.chat_id
                WHERE cm.user_id = $1
                GROUP BY c.id
                LIMIT 25 OFFSET $2
              )
              SELECT 
                uc.chat_id,
                array_agg(ui.first_name || ' ' || ui.second_name) AS members,
                m.id AS message_id,
                m.body AS message_body,
                m.author_id AS message_author_id,
                m.chat_id AS message_chat_id,
                m.created_at AS message_created_at,
                m.readed AS message_readed,
                m.edited AS message_edited,
                m.reply_to AS message_reply_to
              FROM 
                user_chats uc
              LEFT JOIN 
                chat_members cm ON uc.chat_id = cm.chat_id
              LEFT JOIN 
                users u ON cm.user_id = u.id
              LEFT JOIN 
                messages m ON uc.chat_id = m.chat_id
			  INNER JOIN 
			    users_info ui on ui.user_id = u.id
              GROUP BY 
                uc.chat_id, m.id
              ORDER BY 
                uc.chat_id, m.created_at
              LIMIT 50`

	rows, err := r.db.Queryx(query, userId, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chatsMap := make(map[uuid.UUID]*models.Chat)

	for rows.Next() {
		var chatId uuid.UUID
		var members []string
		var message models.Message

		err := rows.Scan(&chatId, pq.Array(&members), &message.Id, &message.Body, &message.AuthorId, &message.ChatId, &message.CreatedAt, &message.Readed, &message.Edited, &message.ReplyTo)
		if err != nil {
			return nil, err
		}

		if message.Body == nil {
			continue
		}

		if chat, exists := chatsMap[chatId]; exists {
			chat.Messages = append(chat.Messages, message)
		} else {
			chatsMap[chatId] = &models.Chat{
				Id:       chatId,
				Members:  members,
				Messages: []models.Message{message},
			}
		}
	}

	var chats []models.Chat
	for _, chat := range chatsMap {
		chats = append(chats, *chat)
	}

	return chats, nil
}

func (r *MessagesPostgres) GetChat(chatId, userId uuid.UUID, page int) (models.Chat, error) {
	// Определяем смещение для пагинации
	offset := (page - 1) * 100

	query := `SELECT 
                c.id AS chat_id,
                array_agg(ui.first_name || ' ' || ui.second_name) AS members,
                m.id AS message_id,
                m.body AS message_body,
                m.author_id AS message_author_id,
                m.chat_id AS message_chat_id,
                m.created_at AS message_created_at,
                m.readed AS message_readed,
                m.edited AS message_edited,
                m.reply_to AS message_reply_to
              FROM 
                chats c
              LEFT JOIN 
                chat_members cm ON c.id = cm.chat_id
              LEFT JOIN 
                users u ON cm.user_id = u.id
              LEFT JOIN 
                messages m ON c.id = m.chat_id
				INNER JOIN 
			    users_info ui on ui.user_id = u.id
              WHERE 
                c.id = $1
              GROUP BY 
                c.id, m.id
              ORDER BY 
                m.created_at
              LIMIT 100 OFFSET $2;`

	rows, err := r.db.Queryx(query, chatId, offset)
	if err != nil {
		return models.Chat{}, err
	}
	defer rows.Close()

	var chat models.Chat
	chat.Id = chatId

	// Сбор участников и сообщений
	for rows.Next() {
		var message models.Message

		err := rows.Scan(&chat.Id, pq.Array(&chat.Members), &message.Id, &message.Body, &message.AuthorId, &message.ChatId, &message.CreatedAt, &message.Readed, &message.Edited, &message.ReplyTo)
		if err != nil {
			return models.Chat{}, err
		}

		// Добавляем сообщение в чат
		chat.Messages = append(chat.Messages, message)
	}

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
