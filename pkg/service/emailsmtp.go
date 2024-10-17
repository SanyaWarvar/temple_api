package service

import (
	"errors"
	"time"

	"github.com/SanyaWarvar/temple_api/pkg/repository"
	"github.com/sirupsen/logrus"
)

type EmailSmtpService struct {
	cache repository.ICacheRepo
	repo  repository.IEmailSmtpRepo
}

func NewEmailSmtpService(repo repository.IEmailSmtpRepo, cache repository.ICacheRepo) *EmailSmtpService {
	return &EmailSmtpService{
		cache: cache,
		repo:  repo,
	}
}

func (s *EmailSmtpService) SendMessage(email, messageText, title string) error {
	return s.repo.SendMessage(email, messageText, title)
}

func (s *EmailSmtpService) SendConfirmEmailMessage(email string) error {
	code := s.GenerateConfirmCode()
	s.cache.SaveConfirmCode(email, code)
	err := s.repo.SendConfirmEmailMessage(email, code)
	if err != nil {
		logrus.Errorf("error while sending confirm email message: %s", err.Error())
	}

	return err
}

func (s *EmailSmtpService) CheckEmailConfirm(email string) (bool, error) {
	return s.repo.CheckEmailConfirm(email)
}

func (s *EmailSmtpService) ConfirmEmail(email, code string) error {
	trueCode, _, err := s.cache.GetConfirmCode(email)
	if err != nil {
		return err
	}
	if trueCode != code {
		return errors.New("bad code")
	}
	return s.repo.ConfirmEmail(email)
}

func (s *EmailSmtpService) GetConfirmCode(email string) (string, time.Duration, error) {
	return s.cache.GetConfirmCode(email)
}

func (s *EmailSmtpService) GenerateConfirmCode() string {
	return s.repo.GenerateConfirmCode()
}
