package service

import (
	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/SanyaWarvar/temple_api/pkg/repository"
	"github.com/google/uuid"
)

type TiktokService struct {
	repo repository.ITiktokRepo
}

func NewTiktokService(repo repository.ITiktokRepo) *TiktokService {
	return &TiktokService{repo: repo}
}

func (s *TiktokService) CreateTiktok(item models.Tiktok) error {
	return s.repo.CreateTiktok(item)
}

func (s *TiktokService) GetTiktokById(tiktokId uuid.UUID) (models.Tiktok, error) {
	return s.repo.GetTiktokById(tiktokId)
}

func (s *TiktokService) DeleteTiktokById(tiktokId, userId uuid.UUID ) error {
	return s.repo.DeleteTiktokById(tiktokId, userId)
}
