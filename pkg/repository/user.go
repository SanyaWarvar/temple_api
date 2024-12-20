package repository

import (
	"fmt"
	"reflect"
	"strings"

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
	tx, err := r.db.DB.Begin()
	if err != nil {
		return err
	}
	id := uuid.NewString()
	query := fmt.Sprintf(`INSERT INTO %s (id, username, email, password_hash) VALUES ($1, $2, $3, $4)`, usersTable)
	_, err = tx.Exec(query, id, user.Username, user.Email, user.Password)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = fmt.Sprintf(`INSERT INTO %s (user_id, first_name) VALUES ($1, $2)`, usersInfoTable)
	_, err = tx.Exec(query, id, user.Username)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
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

func (r *UserPostgres) GetUserByE(email string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf(`SELECT id, username, email, password_hash FROM %s WHERE email = $1`, usersTable)
	err := r.db.Get(&user, query, email)
	return user, err
}

func (r *UserPostgres) GetUserByEP(email, hashedPassword string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf(`SELECT id FROM %s WHERE email = $1 AND password_hash = $2`, usersTable)
	err := r.db.Get(&user, query, email, hashedPassword)
	return user, err
}

func (r *UserPostgres) GetUserById(userId uuid.UUID) (models.User, error) {
	var user models.User
	query := fmt.Sprintf(`SELECT id, username, email, password_hash FROM %s WHERE id = $1`, usersTable)
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

func (r *UserPostgres) GetUserInfoById(userId uuid.UUID) (models.UserInfo, error) {
	var userInfo models.UserInfo
	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1`, usersInfoTable)
	err := r.db.Get(&userInfo, query, userId)
	return userInfo, err
}

func (r *UserPostgres) GetUserInfoByU(username string) (models.UserInfo, error) {
	var userInfo models.UserInfo
	query := fmt.Sprintf(`SELECT user_id, (select username from users where id = ui.user_id) as profile_picture, first_name, second_name, status, birthday, gender, country, city FROM %s ui WHERE user_id = (SELECT id FROM %s WHERE username = $1)`, usersInfoTable, usersTable)
	err := r.db.Get(&userInfo, query, username)
	return userInfo, err
}

type UserInfoWithFS struct {
	models.UserInfo
	FriendStatus string `json:"friend_status" db:"friend_status"`
}

func (r *UserPostgres) GetUserInfoByUWithFS(username string, userId uuid.UUID) (UserInfoWithFS, error) {
	var userInfo UserInfoWithFS
	query := fmt.Sprintf(
		`
		SELECT 
		user_id, 
		(select username from users where id = ui.user_id) as profile_picture, 
		first_name, second_name, status, birthday, gender, country, city,
		coalesce((select case 
			when fi.from_user_id = $1 and fi.confirmed = 't' then 'friends'
			when fi.from_user_id = $1 and fi.confirmed = 'f' then 'follow'
			when fi.from_user_id = ui.user_id and fi.confirmed = 'f' then 'sub'
			else 'not friend'
		end 
		from friends_invites fi
		where (fi.from_user_id = $1 and fi.to_user_id = ui.user_id) 
		or (fi.to_user_id = $1 and fi.from_user_id = ui.user_id)
		), 'not friends') as friend_status
        FROM users_info ui 
        WHERE ui.user_id = (SELECT id FROM users WHERE username = $2) 
		`)
	err := r.db.Get(&userInfo, query, userId, username)
	return userInfo, err
}

func (r *UserPostgres) UpdateUserInfo(userInfo models.UserInfo) error {
	fields := make([]string, 0)
	values := make([]interface{}, 0)

	rv := reflect.ValueOf(userInfo)
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

	query := fmt.Sprintf("UPDATE %s SET %s WHERE user_id = $%d", usersInfoTable, strings.Join(fields, ", "), trueCounter)

	values = append(values, userInfo.UserId)

	_, err := r.db.Exec(query, values...)
	return err
}

type FindUserOutput struct {
	FirstName  string  `json:"first_name" db:"first_name"`
	SecondName string  `json:"second_name" db:"second_name"`
	Username   string  `json:"username" db:"username"`
	ProfilePic string  `json:"profile_picture" db:"profile_picture"`
	Dist       float64 `json:"-" db:"dist"`
}

func (r *UserPostgres) FindUsers(searchString string, page int) ([]FindUserOutput, error) {
	var userInfo []FindUserOutput

	offset := (page - 1) * 50
	var query string
	if strings.HasPrefix(searchString, "@") {
		searchString = searchString[1:]
		query = fmt.Sprintf(`
		SELECT first_name, second_name, (select username from users where id = ui.user_id) as profile_picture, (select username from users where id = ui.user_id) FROM %s ui where user_id = (SELECT id FROM %s WHERE username = $1)
		LIMIT 50
		OFFSET $2
		 `, usersInfoTable, usersTable,
		)
	} else {
		query = fmt.Sprintf(`
		SELECT first_name, second_name, fullname(first_name, second_name) <-> $1 as dist, (select username from users where id = ui.user_id) as profile_picture, (select username from users where id = ui.user_id)
		from %s ui order by dist
		LIMIT 50
		OFFSET $2
		 `, usersInfoTable,
		)
	}

	err := r.db.Select(&userInfo, query, searchString, offset)
	return userInfo, err
}

func (r *UserPostgres) UpdateProfPic(userId uuid.UUID, path string) error {

	query := fmt.Sprintf(
		`
		UPDATE %s
		SET profile_picture = $1
		WHERE user_id = $2
		`,
		usersInfoTable,
	)

	_, err := r.db.Exec(query, path, userId)
	return err
}
