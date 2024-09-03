package service

import (
	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/SanyaWarvar/temple_api/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) error
	GetUser(user models.User) (models.User, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
