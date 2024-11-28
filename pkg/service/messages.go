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

func (r *MessagesService) CreateChat(inviteUsername string, owner uuid.UUID) (uuid.UUID, error) {
	return r.repo.CreateChat(inviteUsername, owner)
}

func (r *MessagesService) GetAllChats(userId uuid.UUID, page int) ([]models.Chat, error) {
	return r.repo.GetAllChats(userId, page)
}

func (r *MessagesService) GetChat(chatId, userId uuid.UUID, page int) (models.Chat, error) {
	return r.repo.GetChat(chatId, userId, page)
}

func (r *MessagesService) CreateMessage(data models.Message) error {
	return r.repo.CreateMessage(data)
}

func (r *MessagesService) ReadMessage(messageId, userId uuid.UUID) error {
	return r.repo.ReadMessage(messageId, userId)
}

func (r *MessagesService) EditMessage(userId uuid.UUID, message models.Message) error {
	return r.repo.EditMessage(userId, message)
}

func (r *MessagesService) DeleteMessage(messageId, userId uuid.UUID) error {
	return r.repo.DeleteMessage(messageId, userId)
}
