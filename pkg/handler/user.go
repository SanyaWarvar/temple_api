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
	userId, err := getUserId(c)
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
