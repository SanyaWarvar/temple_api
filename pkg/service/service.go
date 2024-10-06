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
	HashPassword(password string) (string, error)
	//нужно ли мне это???
	//GetUserById(userId uuid.UUID) (models.User, error)
	//ComparePassword(password, hashedPassword string) bool
}

type IEmailSmtpService interface {
	CheckEmailConfirm(email string) (bool, error)
	ConfirmEmail(email, code string) error
	SendConfirmEmailMessage(user models.User) error
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

type Service struct {
	IUserService
	IEmailSmtpService
	IJwtManagerService
	ICacheService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		IUserService:       NewUserService(repos.IUserRepo),
		IEmailSmtpService:  NewEmailSmtpService(repos.IEmailSmtpRepo, repos.ICacheRepo),
		IJwtManagerService: NewJwtManagerService(repos.IJwtManagerRepo),
		ICacheService:      NewCacheService(repos.ICacheRepo),
	}
}
