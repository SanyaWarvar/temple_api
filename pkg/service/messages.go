package service

import (
	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/SanyaWarvar/temple_api/pkg/repository"
	"github.com/google/uuid"
)

type MessagesService struct {
	repo repository.IMessagesRepo
}

func NewMessagesService(repo repository.IMessagesRepo) *MessagesService {
	return &MessagesService{repo: repo}
}

func (s *MessagesService) CreateChat(inviteUsername string, owner uuid.UUID) (uuid.UUID, error) {
	return s.repo.CreateChat(inviteUsername, owner)
}

func (s *MessagesService) GetAllChats(userId uuid.UUID, page int) ([]repository.AllChatsOutput, error) {
	return s.repo.GetAllChats(userId, page)
}

func (s *MessagesService) GetChat(chatId, userId uuid.UUID, page int) (repository.AllChatsOutput, error) {
	return s.repo.GetChat(chatId, userId, page)
}

func (s *MessagesService) CreateMessage(data models.Message) error {
	return s.repo.CreateMessage(data)
}

func (s *MessagesService) ReadMessage(messageId, userId uuid.UUID) error {
	return s.repo.ReadMessage(messageId, userId)
}

func (s *MessagesService) EditMessage(userId uuid.UUID, message models.Message) error {
	return s.repo.EditMessage(userId, message)
}

func (s *MessagesService) DeleteMessage(messageId, userId uuid.UUID) error {
	return s.repo.DeleteMessage(messageId, userId)
}

func (s *MessagesService) GetMembersFromChatByID(chatId uuid.UUID) ([]models.User, error) {
	return s.repo.GetMembersFromChatByID(chatId)
}
