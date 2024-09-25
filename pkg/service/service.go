package service

import (
	"time"

	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/SanyaWarvar/temple_api/pkg/repository"
)

type Authorizer interface {
	CreateUser(user models.User) error
	GetUser(user models.User) (models.User, error)
	CheckEmailConfirm(email string) (bool, error)
	ConfirmEmail(email, code string) error
	GetConfirmCode(email string) (string, time.Duration, error)
}

type EmailSmtper interface {
	SendMessage(email, messageText, title string) error
	SendConfirmEmailMessage(user models.User) error
}

type Service struct {
	Authorizer
	EmailSmtper
}

func NewService(repos *repository.Repository, authSettings authSettings, emailSmtpSettings EmailSettings) *Service {
	return &Service{
		Authorizer:  NewAuthService(repos.Authorizer, repos.Cacher, authSettings),
		EmailSmtper: NewEmailSmtpService(repos.Cacher, emailSmtpSettings),
	}
}
