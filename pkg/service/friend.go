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

func (s *FriendService) GetAllFriend(username string, page int) (repository.FriendListOutput, error) {
	return s.repo.GetAllFriend(username, page)
}

func (s *FriendService) GetAllSubs(username string, page int) (repository.SubListOutput, error) {
	return s.repo.GetAllSubs(username, page)
}

func (s *FriendService) GetAllFollows(username string, page int) (repository.FollowListOutput, error) {
	return s.repo.GetAllFollows(username, page)
}

func (s *FriendService) CheckFriendStatus(fromId, toId uuid.UUID) (repository.FriendStatus, error) {
	return s.repo.CheckFriendStatus(fromId, toId)
}
