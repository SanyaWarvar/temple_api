package handler

import (
	"net/http"

	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) getUserInfo(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		newErrorResponse(c, http.StatusBadRequest, "bad username")
		return
	}

	userInfo, err := h.services.IUserService.GetUserInfoByU(username)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, userInfo)
}

func (h *Handler) updateUserInfo(c *gin.Context) {
	userId, err := getUserId(c, true)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Can't parse user id")
		return
	}

	var input models.UserInfo

	err = c.BindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "bad json")
		return
	}
	input.UserId = userId

	err = h.services.IUserService.UpdateUserInfo(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

type findUserStruct struct {
	SearchString string `json:"search_string" binding:"required"`
	Page         int    `json:"page"`
}

func (h *Handler) findUser(c *gin.Context) {
	var input findUserStruct
	c.BindJSON(&input)
	users, err := h.services.IUserService.FindUsers(input.SearchString, input.Page)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users)
}
