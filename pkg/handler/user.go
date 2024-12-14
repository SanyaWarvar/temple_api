package handler

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"

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

	file, err := os.OpenFile("user_data/profile_pictures/"+*userInfo.ProfilePic, os.O_RDONLY, 0666)
	if err != nil {
		temp := c.Request.Host + "/images/base/base_pic.jpg"
		userInfo.ProfilePic = &temp
	} else {
		temp := c.Request.Host + "/images/profiles/" + *userInfo.ProfilePic
		userInfo.ProfilePic = &temp
		file.Close()
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

	for ind, item := range users {
		file, err := os.OpenFile("user_data/profile_pictures/"+item.ProfilePic, os.O_RDONLY, 0666)
		if err != nil {
			users[ind].ProfilePic = c.Request.Host + "/images/base/base_pic.jpg"
		} else {
			users[ind].ProfilePic = c.Request.Host + "/images/profiles/" + item.ProfilePic
			file.Close()
		}
	}

	c.JSON(http.StatusOK, users)
}

type UpdateProfPicInput struct {
	ProfilePic *multipart.FileHeader `form:"profile_pic" binding:"required"`
}

func (h *Handler) updateProfPic(c *gin.Context) {
	var input UpdateProfPicInput
	userId, _ := getUserId(c, false)
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
	suffix := filepath.Ext(input.ProfilePic.Filename)
	ValidFileSuffixForProfilePicture := []string{".gif", ".jpg", ".png", ".svg"}
	if !slices.Contains(ValidFileSuffixForProfilePicture, suffix) {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Supported formats: .gif, .jpg, .png, .svg. %s is unsupported!", suffix))
		return
	}

	path := base64.RawStdEncoding.EncodeToString(fileBytes)

	err = h.services.IUserService.UpdateProfPic(userId, path)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	os.WriteFile(path, fileBytes, 0644)

	c.JSON(http.StatusCreated, map[string]string{"details": "success"})
}
