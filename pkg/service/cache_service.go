package service

import (
	"time"

	"github.com/SanyaWarvar/temple_api/pkg/repository"
)

type CacheService struct {
	cache repository.ICacheRepo
}

func NewCacheService(cache repository.ICacheRepo) *CacheService {
	return &CacheService{cache: cache}
}

func (s *CacheService) GetConfirmCode(email string) (string, time.Duration, error) {
	return s.cache.GetConfirmCode(email)
}
func (s *CacheService) SaveConfirmCode(email, code string) error {
	return s.cache.SaveConfirmCode(email, code)
}
