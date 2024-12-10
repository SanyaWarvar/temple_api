package repository

import (
	"time"

	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type IUserRepo interface {
	CreateUser(user models.User) error
	GetUserByUP(username, hashedPassword string) (models.User, error)
	GetUserByEP(email, hashedPassword string) (models.User, error)
	GetUserById(userId uuid.UUID) (models.User, error)
	GetUserByU(username string) (models.User, error)
	GetUserByE(email string) (models.User, error)
	HashPassword(password string) (string, error)
	ComparePassword(password, hashedPassword string) bool
	GetUserInfoById(userId uuid.UUID) (models.UserInfo, error)
	UpdateUserInfo(userInfo models.UserInfo) error
	GetUserInfoByU(username string) (models.UserInfo, error)
	FindUsers(searchString string, page int) ([]FindUserOutput, error)
	UpdateProfPic(userId uuid.UUID, path string) error
}

type IEmailSmtpRepo interface {
	CheckEmailConfirm(email string) (bool, error)
	ConfirmEmail(email string) error
	SendConfirmEmailMessage(email, code string) error
	SendMessage(email, messageText, title string) error
	GenerateConfirmCode() string
}

type IJwtManagerRepo interface {
	GenerateAccessToken(userId, refreshId uuid.UUID) (string, error)
	GenerateRefreshToken(userId uuid.UUID) (string, error)
	GeneratePairToken(userId uuid.UUID) (string, string, uuid.UUID, error)
	CompareTokens(hashedToken, token string) bool
	HashToken(refreshToken string) (string, error)
	SaveRefreshToken(hashedToken string, tokenId, userId uuid.UUID) error
	DeleteRefreshTokenById(tokenId uuid.UUID) error
	GetRefreshTokenById(tokenId uuid.UUID) (string, error)
	ParseToken(accessToken string) (*models.AccessTokenClaims, error)
	CheckRefreshTokenExp(tokenId uuid.UUID) bool
}

type ICacheRepo interface {
	//emailsmtp
	GetConfirmCode(email string) (string, time.Duration, error)
	SaveConfirmCode(email, code string) error

	//friends
}

type IFriendRepo interface {
	InviteFriend(fromId uuid.UUID, toUsername string) error
	DeleteByU(invitedId uuid.UUID, ownerUsername string) error
	ConfirmFriend(invitedId uuid.UUID, ownerUsername string) error
	GetAllFriend(userId uuid.UUID, page int) (FriendListOutput, error)
}

type IUsersPostsRepo interface {
	CreatePost(post models.UserPost) error
	UpdatePost(newPost models.UserPost) error
	GetPostById(postId, userId uuid.UUID) (UserPostOutput, error)
	DeletePostById(postId, userId uuid.UUID) error

	GetPostsByU(username string, page int, userId uuid.UUID) ([]UserPostOutput, error)

	LikePostById(postId, userId uuid.UUID) error
}

type IMessagesRepo interface {
	CreateChat(inviteUsername string, owner uuid.UUID) (uuid.UUID, error)
	GetAllChats(userId uuid.UUID, page int) ([]models.Chat, error)
	GetChat(chatId, userId uuid.UUID, page int) (models.Chat, error)

	CreateMessage(data models.Message) error
	ReadMessage(messageId, userId uuid.UUID) error
	EditMessage(userId uuid.UUID, message models.Message) error
	DeleteMessage(messageId, userId uuid.UUID) error
}

type ITiktokRepo interface {
	CreateTiktok(item models.Tiktok) error
	GetTiktokById(tiktokId uuid.UUID) (models.Tiktok, error)
	DeleteTiktokById(tiktokId, userId uuid.UUID) error
}

type Repository struct {
	IUserRepo
	IEmailSmtpRepo
	IJwtManagerRepo
	ICacheRepo
	IFriendRepo
	IUsersPostsRepo
	IMessagesRepo
	ITiktokRepo
}

func NewRepository(db *sqlx.DB, cacheDb *redis.Client, codeExp time.Duration, emailCfg *EmailCfg, jwtCfg *JwtManagerCfg) *Repository {
	return &Repository{
		IUserRepo:       NewUserPostgres(db),
		IEmailSmtpRepo:  NewEmailSmtpPostgres(db, emailCfg),
		IJwtManagerRepo: NewJwtManagerPostgres(db, jwtCfg),
		ICacheRepo:      NewCacheRedis(cacheDb, codeExp),
		IFriendRepo:     NewFriendPostgres(db),
		IUsersPostsRepo: NewUsersPostsPostgres(db),
		IMessagesRepo:   NewMessagesPostgres(db),
		ITiktokRepo:     NewTiktokPostgres(db),
	}
}
