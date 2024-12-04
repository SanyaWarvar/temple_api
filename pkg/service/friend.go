package service

import (
	"github.com/SanyaWarvar/temple_api/pkg/repository"
	"github.com/google/uuid"
)

type FriendService struct {
	repo repository.IFriendRepo
}

func NewFriendService(repo repository.IFriendRepo) *FriendService {
	return &FriendService{repo: repo}
}

func (s *FriendService) InviteFriend(fromId uuid.UUID, toUsername string) error {
	return s.repo.InviteFriend(fromId, toUsername)
}

func (s *FriendService) DeleteByU(invitedId uuid.UUID, ownerUsername string) error {
	return s.repo.DeleteByU(invitedId, ownerUsername)
}

func (s *FriendService) ConfirmFriend(invitedId uuid.UUID, ownerUsername string) error {
	return s.repo.ConfirmFriend(invitedId, ownerUsername)
}

func (s *FriendService) GetAllFriend(userId uuid.UUID, page int) (repository.FriendListOutput, error) {
	return s.repo.GetAllFriend(userId, page)
}
