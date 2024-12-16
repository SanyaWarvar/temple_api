package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/SanyaWarvar/temple_api/pkg/repository"
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

	go h.PrepareFriendNotify(c, requestOwnerId, toUsername, repository.Follow)

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

type UsernameInput struct {
	Username string `json:"username"`
}

func (h *Handler) getAllFriends(c *gin.Context) {
	var input UsernameInput
	err := c.BindJSON(&input)
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
	friends, err := h.services.IFriendService.GetAllFriend(input.Username, pageInt)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, friends)
}

func (h *Handler) getAllSubs(c *gin.Context) {
	var input UsernameInput
	err := c.BindJSON(&input)
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
	subs, err := h.services.IFriendService.GetAllSubs(input.Username, pageInt)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, subs)
}

func (h *Handler) getAllFollows(c *gin.Context) {
	var input UsernameInput
	err := c.BindJSON(&input)
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
	follows, err := h.services.IFriendService.GetAllFollows(input.Username, pageInt)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, follows)
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

	go h.PrepareFriendNotify(c, OwnerId, Username, repository.Confirmed)
	c.JSON(http.StatusNoContent, nil)
}
