package handler

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/SanyaWarvar/temple_api/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Notify struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type NewMessageNotify2 struct {
	AuthorUsername string `json:"from_username"`
	Message        string `json:"message"`
}

func (h *Handler) NewMessageNotify(authorUsername, recipientUsername, message string) {
	msg := Notify{
		Type: "new_message",
		Data: NewMessageNotify2{
			AuthorUsername: authorUsername,
			Message:        message,
		},
	}
	targetConn, ok := h.clients[recipientUsername]
	if !ok {
		return
	}
	data, err := json.Marshal(msg)
	fmt.Println(err)

	targetConn.WriteMessage(1, data)
}

type FriendNotify struct {
	Username   string                  `json:"username"`
	FirstName  string                  `json:"first_name"`
	SecondName string                  `json:"second_name"`
	ProfilePic string                  `json:"profile_picture"`
	Status     repository.FriendStatus `json:"status"`
}

func (h *Handler) PrepareFriendNotify(c *gin.Context, requestOwnerId uuid.UUID, toUsername string, status repository.FriendStatus) {

	targetUserInfo, err := h.services.IUserService.GetUserInfoByU(toUsername)
	if err != nil {
		logrus.Errorf("Error while send friend notify: %s", err.Error())
		return
	}
	targetUser, err := h.services.IUserService.GetUserById(targetUserInfo.UserId)
	if err != nil {
		logrus.Errorf("Error while send friend notify: %s", err.Error())
		return
	}

	fromUser, err := h.services.IUserService.GetUserById(requestOwnerId)
	if err != nil {
		logrus.Errorf("Error while send friend notify: %s", err.Error())
		return
	}
	fromUserInfo, err := h.services.IUserService.GetUserInfoById(requestOwnerId)
	if err != nil {
		logrus.Errorf("Error while send friend notify: %s", err.Error())
		return
	}

	//TODO вынести эту логику
	file, err := os.OpenFile("user_data/profile_pictures/"+*fromUserInfo.ProfilePic, os.O_RDONLY, 0666)
	if err != nil {
		temp := c.Request.Host + "/images/base/base_pic.jpg"
		fromUserInfo.ProfilePic = &temp
	} else {
		temp := c.Request.Host + "/images/profiles/" + *fromUserInfo.ProfilePic
		fromUserInfo.ProfilePic = &temp
		file.Close()
	}

	h.FriendNotify(targetUser.Username, fromUser.Username, status, fromUserInfo)

}

func (h *Handler) FriendNotify(recipientUsername, fromUsername string, status repository.FriendStatus, fromUserInfo models.UserInfo) {
	msg := Notify{
		Type: "friend",
		Data: FriendNotify{
			Username:   fromUsername,
			FirstName:  *fromUserInfo.FirstName,
			SecondName: *fromUserInfo.SecondName,
			ProfilePic: *fromUserInfo.ProfilePic,
			Status:     status,
		},
	}
	targetConn, ok := h.clients[recipientUsername]
	if !ok {
		return
	}

	targetConn.WriteJSON(msg)
}
