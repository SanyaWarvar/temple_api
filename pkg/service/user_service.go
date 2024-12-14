package service

import (
	"errors"

	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/SanyaWarvar/temple_api/pkg/repository"
	"github.com/google/uuid"
)

type UserService struct {
	repo repository.IUserRepo
}

func NewUserService(repo repository.IUserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user models.User) error {
	var err error
	user.Password, err = s.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Id, err = uuid.NewUUID()
	if err != nil {
		return err
	}
	return s.repo.CreateUser(user)
}

func (s *UserService) GetUserByUP(user models.User) (models.User, error) {
	targetUser, err := s.repo.GetUserByU(user.Username)
	if err != nil {
		return user, err
	}

	if s.repo.ComparePassword(user.Password, targetUser.Password) {
		return targetUser, err
	}
	return user, errors.New("incorrect password")
}

func (s *UserService) GetUserByEP(email, password string) (models.User, error) {
	var user models.User
	targetUser, err := s.repo.GetUserByE(email)
	if err != nil {
		return user, err
	}

	if s.repo.ComparePassword(password, targetUser.Password) {
		return targetUser, err
	}
	return user, errors.New("incorrect password")
}

func (s *UserService) HashPassword(password string) (string, error) {
	return s.repo.HashPassword(password)
}

func (s *UserService) GetUserInfoById(userId uuid.UUID) (models.UserInfo, error) {
	return s.repo.GetUserInfoById(userId)
}
func (s *UserService) UpdateUserInfo(userInfo models.UserInfo) error {
	return s.repo.UpdateUserInfo(userInfo)
}

func (s *UserService) GetUserInfoByU(username string) (models.UserInfo, error) {
	return s.repo.GetUserInfoByU(username)
}

func (s *UserService) FindUsers(searchString string, page int) ([]repository.FindUserOutput, error) {
	return s.repo.FindUsers(searchString, page)
}

func (s *UserService) UpdateProfPic(userId uuid.UUID, path string) error {
	return s.repo.UpdateProfPic(userId, path)
}

func (s *UserService) GetUserById(userId uuid.UUID) (models.User, error) {
	return s.repo.GetUserById(userId)
}
