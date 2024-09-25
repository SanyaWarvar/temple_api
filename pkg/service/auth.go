package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/SanyaWarvar/temple_api/pkg/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}

type authSettings struct {
	TokenTTL   time.Duration
	Salt       string
	SigningKey string
}

func NewAuthSettings(tokenTTL time.Duration, salt, signingKey string) *authSettings {
	return &authSettings{
		TokenTTL:   tokenTTL,
		Salt:       salt,
		SigningKey: signingKey,
	}
}

type AuthService struct {
	repo     repository.Authorizer
	cache    repository.Cacher
	settings authSettings
}

func NewAuthService(repo repository.Authorizer, cache repository.Cacher, settings authSettings) *AuthService {
	return &AuthService{repo: repo, cache: cache, settings: settings}
}

func (s *AuthService) CreateUser(user models.User) error {
	user.Password = s.generatePasswordHash(user.Password)
	user.Id = uuid.NewString()
	return s.repo.CreateUser(user)
}

func (s *AuthService) GetUser(user models.User) (models.User, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.GetUser(user.Username, user.Password)
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(s.settings.Salt)))
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&tokenClaims{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(s.settings.TokenTTL).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			user.Id,
		},
	)

	return token.SignedString([]byte(s.settings.SigningKey))
}

func (s *AuthService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.settings.SigningKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are incorrect")
	}

	return claims.UserId, nil
}

func (s *AuthService) CheckEmailConfirm(email string) (bool, error) {

	return s.repo.CheckEmailConfirm(email)

}

func (s *AuthService) ConfirmEmail(email, code string) error {
	trueCode, _, err := s.cache.GetConfirmCode(email)
	if err != nil {
		return err
	}
	if trueCode != code {
		return errors.New("bad code")
	}
	return s.repo.ConfirmEmail(email)
}

func (s *AuthService) GetConfirmCode(email string) (string, time.Duration, error) {

	return s.cache.GetConfirmCode(email)
}
