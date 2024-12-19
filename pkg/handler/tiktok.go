package handler

import (
	"encoding/base64"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateTiktokInput struct {
	File  *multipart.FileHeader `form:"file" binding:"required"`
	Title string                `form:"title" binding:"required"`
}

func (h *Handler) createTiktok(c *gin.Context) {
	var input CreateTiktokInput
	userId, _ := getUserId(c, false)

	if err := c.ShouldBind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	file, err := input.File.Open()
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	size := input.File.Size / 1024 / 1024
	if size > 50 { //todo вынести 50 в конфиг
		newErrorResponse(c, http.StatusBadRequest, "max size is 50mb")
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	body := base64.RawStdEncoding.EncodeToString(fileBytes)
	ext := filepath.Ext(input.File.Filename)
	item := models.Tiktok{
		Id:        uuid.New(),
		AuthorId:  userId,
		CreatedAt: time.Now(),
		Title:     input.Title + ext,
		Body:      body,
	}

	err = h.services.ITiktokService.CreateTiktok(item)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	go func() {
		os.WriteFile("user_data/tik_toks/"+item.Id.String(), fileBytes, 0644)
	}()

	c.JSON(http.StatusCreated, map[string]string{"id": item.Id.String()})
}

func (h *Handler) getTiktokById(c *gin.Context) {
	tiktokIdStr := c.Param("id")
	tiktokId, err := uuid.Parse(tiktokIdStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	item, err := h.services.ITiktokService.GetTiktokById(tiktokId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	item.Body = c.Request.Host + "/tik_toks/" + item.Id.String()
	c.JSON(http.StatusOK, item)
}

func (h *Handler) deleteTiktokById(c *gin.Context) {
	tiktokIdStr := c.Param("id")
	tiktokId, err := uuid.Parse(tiktokIdStr)
	userId, _ := getUserId(c, false)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.ITiktokService.DeleteTiktokById(tiktokId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) feedTikTok(c *gin.Context) {
	userId, err := getUserId(c, false)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var input PageInput
	err = c.BindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	output, err := h.services.ITiktokService.Feed(userId, input.Page)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	for ind, item := range output {
		file, err := os.OpenFile("user_data/profile_pictures/"+item.AuthorProfilePicture, os.O_RDONLY, 0666)
		if err != nil {
			output[ind].AuthorProfilePicture = c.Request.Host + "/images/base/base_pic.jpg"
		} else {
			output[ind].AuthorProfilePicture = c.Request.Host + "/images/profiles/" + item.AuthorProfilePicture
			file.Close()
		}

		file, err = os.OpenFile("user_data/tik_toks/"+item.Id.String(), os.O_RDONLY, 0666)
		if err != nil {
			bytesData, err := base64.RawStdEncoding.DecodeString(item.Body)
			if err != nil {
				continue
			}
			os.WriteFile("user_data/tik_toks/"+item.Id.String(), bytesData, 0644)
		} else {
			output[ind].Body = c.Request.Host + "/tik_toks/" + item.Id.String()
			file.Close()
		}
		output[ind].Body = c.Request.Host + "/tik_toks/" + item.Id.String()
	}

	c.JSON(http.StatusOK, output)
}
