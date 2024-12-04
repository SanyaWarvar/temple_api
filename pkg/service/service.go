package service

import (
	"time"

	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/SanyaWarvar/temple_api/pkg/repository"
	"github.com/google/uuid"
)

type IUserService interface {
	CreateUser(user models.User) error
	GetUserByUP(user models.User) (models.User, error)
	GetUserByEP(email, password string) (models.User, error)
	HashPassword(password string) (string, error)
	//нужно ли мне это???
	//GetUserById(userId uuid.UUID) (models.User, error)
	//ComparePassword(password, hashedPassword string) bool
	GetUserInfoById(userId uuid.UUID) (models.UserInfo, error)
	GetUserInfoByU(username string) (models.UserInfo, error)
	UpdateUserInfo(userInfo models.UserInfo) error
	FindUsers(searchString string, page int) ([]repository.FindUserOutput, error)
	UpdateProfPic(userId uuid.UUID, path string) error
}

type IEmailSmtpService interface {
	CheckEmailConfirm(email string) (bool, error)
	ConfirmEmail(email, code string) error
	SendConfirmEmailMessage(email string) error
	SendMessage(email, messageText, title string) error
	GenerateConfirmCode() string
}

type IJwtManagerService interface {
	/*Будто бы я и так всегда парами генерю и не надо по отдельности?
	GenerateAccessToken(userId, refreshId uuid.UUID) (string, error)
	GenerateRefreshToken(userId uuid.UUID) (string, error)
	*/
	GeneratePairToken(userId uuid.UUID) (string, string, uuid.UUID, error)
	CompareTokens(refreshId uuid.UUID, token string) bool
	//HashToken(refreshToken string) (string, error)
	SaveRefreshToken(hashedToken string, userId, tokenId uuid.UUID) error
	DeleteRefreshTokenById(tokenId uuid.UUID) error
	GetRefreshTokenById(tokenId uuid.UUID) (string, error)
	ParseToken(accessToken string) (*models.AccessTokenClaims, error)
	CheckRefreshTokenExp(tokenId uuid.UUID) bool
}

type ICacheService interface {
	GetConfirmCode(email string) (string, time.Duration, error)
	SaveConfirmCode(email, code string) error
}

type IFriendService interface {
	InviteFriend(fromId uuid.UUID, toUsername string) error
	DeleteByU(invitedId uuid.UUID, ownerUsername string) error
	ConfirmFriend(invitedId uuid.UUID, ownerUsername string) error
	GetAllFriend(userId uuid.UUID, page int) (repository.FriendListOutput, error)
}

type IUsersPostsService interface {
	CreatePost(post models.UserPost) (uuid.UUID, error)
	UpdatePost(newPost models.UserPost) error
	GetPostById(postId, userId uuid.UUID) (repository.UserPostOutput, error)
	DeletePostById(postId, userId uuid.UUID) error

	GetPostsByU(username string, page int, userId uuid.UUID) ([]repository.UserPostOutput, error)

	LikePostById(postId, userId uuid.UUID) error
}

type IMessagesService interface {
	CreateChat(inviteUsername string, owner uuid.UUID) (uuid.UUID, error)
	GetAllChats(userId uuid.UUID, page int) ([]models.Chat, error)
	GetChat(chatId, userId uuid.UUID, page int) (models.Chat, error)

	CreateMessage(data models.Message) error
	ReadMessage(messageId, userId uuid.UUID) error
	EditMessage(userId uuid.UUID, message models.Message) error
	DeleteMessage(messageId, userId uuid.UUID) error
}

type Service struct {
	IUserService
	IEmailSmtpService
	IJwtManagerService
	ICacheService
	IFriendService
	IUsersPostsService
	IMessagesService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		IUserService:       NewUserService(repos.IUserRepo),
		IEmailSmtpService:  NewEmailSmtpService(repos.IEmailSmtpRepo, repos.ICacheRepo),
		IJwtManagerService: NewJwtManagerService(repos.IJwtManagerRepo),
		ICacheService:      NewCacheService(repos.ICacheRepo),
		IFriendService:     NewFriendService(repos.IFriendRepo),
		IUsersPostsService: NewUsersPostsService(repos.IUsersPostsRepo),
		IMessagesService:   NewMessagesService(repos.IMessagesRepo),
	}
}
