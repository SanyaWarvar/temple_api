package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) inviteFriend(c *gin.Context) {
	requestOwnerId, err := getUserId(c, true)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	toUsername := c.Param("username")
	if toUsername == "" {
		newErrorResponse(c, http.StatusBadRequest, "bad username")
	}
	err = h.services.IFriendService.InviteFriend(requestOwnerId, toUsername)
	if err != nil {
		if strings.Contains(err.Error(), "pkey") {
			newErrorResponse(c, http.StatusBadRequest, "already invited")
		} else if strings.Contains(err.Error(), "check") {
			newErrorResponse(c, http.StatusBadRequest, "can't invite self")
		} else if strings.Contains(err.Error(), "not-null") {
			newErrorResponse(c, http.StatusBadRequest, "username not found")
		}
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) deleteFriend(c *gin.Context) {
	OwnerId, err := getUserId(c, true)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	Username := c.Param("username")
	err = h.services.IFriendService.DeleteByU(OwnerId, Username)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) getAllFriends(c *gin.Context) {
	userId, err := getUserId(c, false)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	pageStr := c.Param("page")
	pageInt, err := strconv.Atoi(pageStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	friends, err := h.services.IFriendService.GetAllFriend(userId, pageInt)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	for i := range friends.Friends {
		if friends.Friends[i].ProfilePicUrl == "base_media/base_pic.jpg" {
			friends.Friends[i].ProfilePicUrl = (c.Request.Host + "/images/base/" + "base_pic.jpg")
		} else {
			friends.Friends[i].ProfilePicUrl = (c.Request.Host + "/images/profiles" + strings.Replace(friends.Friends[i].ProfilePicUrl, "user_data/profile_pictures", "", 1))

		}

		fmt.Println(friends.Friends[i])
	}

	c.JSON(http.StatusOK, friends)
}

func (h *Handler) confirmFriend(c *gin.Context) {
	OwnerId, err := getUserId(c, true)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	Username := c.Param("username")
	err = h.services.IFriendService.ConfirmFriend(OwnerId, Username)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
