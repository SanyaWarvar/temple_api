package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

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
	if input.Page == 0 {
		input.Page = 1
	}
	users, err := h.services.IUserService.FindUsers(input.SearchString, input.Page)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, users)
}

type UpdateProfPicInput struct {
	ProfilePic *multipart.FileHeader `form:"profile_pic" binding:"required"`
}

func (h *Handler) updateProfPic(c *gin.Context) {
	var input UpdateProfPicInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	file, err := input.ProfilePic.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open file"})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read file"})
		return
	}
	os.Mkdir("user_data", 0755)
	os.WriteFile(fmt.Sprintf("user_data/%s", input.ProfilePic.Filename), fileBytes, 0644)

	c.JSON(http.StatusOK, "")
}

type getProfPicInput struct {
	Path string `json:"path"`
}

func (h *Handler) getProfPic(c *gin.Context) {
	var input getProfPicInput

	c.BindJSON(&input)

	file, err := os.Open(input.Path)
	fmt.Println(err)

	defer file.Close()

	// Читаем файл в байты
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read file"})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", fileBytes)
}
