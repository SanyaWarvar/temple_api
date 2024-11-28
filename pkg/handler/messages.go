package handler

import (
	"fmt"
	"net/http"

	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateChatInput struct {
	Username string `json:"username" binding:"required"`
}

func (h *Handler) CreateChat(c *gin.Context) {
	var input CreateChatInput
	err := c.BindJSON(&input)
	userId, _ := getUserId(c, false)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return

	}
	chatId, err := h.services.CreateChat(input.Username, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]string{"chatId": chatId.String()})
}

type PageInput struct {
	Page int `json:"page" binding:"required"`
}

func (h *Handler) GetAllChats(c *gin.Context) {
	var input PageInput
	err := c.BindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, _ := getUserId(c, false)
	fmt.Println(userId)
	chats, err := h.services.GetAllChats(userId, input.Page)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, chats)
}

func (h *Handler) GetChat(c *gin.Context) {
	var input PageInput
	chat_id_string := c.Param("chat_id")
	chat_id, err := uuid.Parse(chat_id_string)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = c.BindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	userId, _ := getUserId(c, false)
	data, err := h.services.GetChat(chat_id, userId, input.Page)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

type NewMessage struct {
	Body string `json:"body" binding:"required"`
}

func (h *Handler) NewMessage(c *gin.Context) {
	var input NewMessage
	c.BindJSON(&input)
	userId, _ := getUserId(c, false)
	chatIdString := c.Param("chat_id")
	chatId, err := uuid.Parse(chatIdString)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	emptyUUID, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	message := models.Message{
		Id:        uuid.New(),
		Body:      &input.Body,
		AuthorId:  userId,
		ChatId:    chatId,
		CreatedAt: nil,
		Readed:    nil,
		Edited:    nil,
		ReplyTo:   emptyUUID,
	}

	err = h.services.CreateMessage(message)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]string{"message_id": message.Id.String()})
}

func (h *Handler) ReadMessage(c *gin.Context) {
	messageIdString := c.Param("message_id")
	messageId, err := uuid.Parse(messageIdString)
	if err != nil {

	}
	userId, _ := getUserId(c, false)
	err = h.services.IMessagesService.ReadMessage(messageId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
func (h *Handler) EditMessage(c *gin.Context) {
	userId, _ := getUserId(c, false)
	var input models.Message

	c.BindJSON(&input)
	err := h.services.IMessagesService.EditMessage(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
func (h *Handler) DeleteMessage(c *gin.Context) {
	messageIdString := c.Param("message_id")
	messageId, err := uuid.Parse(messageIdString)
	if err != nil {

	}
	userId, _ := getUserId(c, false)

	err = h.services.IMessagesService.DeleteMessage(messageId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
