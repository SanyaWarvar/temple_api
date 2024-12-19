package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) getPost(c *gin.Context) {
	userId, err := getUserId(c, true)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	postId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.services.IUsersPostsService.GetPostById(postId, userId)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *Handler) createPost(c *gin.Context) {
	userId, err := getUserId(c, true)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var post models.UserPost
	c.BindJSON(&post)
	post.AuthorId = userId
	post.LastUpdate = time.Now()

	id, err := h.services.IUsersPostsService.CreatePost(post)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, map[string]interface{}{"post_id": id})
}

func (h *Handler) deletePost(c *gin.Context) {
	userId, err := getUserId(c, true)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	postId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.IUsersPostsService.DeletePostById(postId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusConflict, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"details": "success"})
}

func (h *Handler) updatePost(c *gin.Context) {
	userId, err := getUserId(c, true)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var post models.UserPost
	c.BindJSON(&post)
	postId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	post.Id = postId
	post.AuthorId = userId
	post.LastUpdate = time.Now()
	err = h.services.IUsersPostsService.UpdatePost(post)
	if err != nil {
		newErrorResponse(c, http.StatusConflict, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"details": "success"})
}

type getPostsByUInput struct {
	Page int `json:"page"`
}

func (h *Handler) getPostsByU(c *gin.Context) {
	var input getPostsByUInput
	userId, err := getUserId(c, true)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	username := c.Param("username")
	err = c.BindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(username)
	data, err := h.services.IUsersPostsService.GetPostsByU(username, input.Page, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *Handler) likePost(c *gin.Context) {
	userId, err := getUserId(c, true)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	postId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.IUsersPostsService.LikePostById(postId, userId)
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			newErrorResponse(c, http.StatusConflict, "post not found")
			return
		} else {
			newErrorResponse(c, http.StatusConflict, err.Error())
			return
		}

	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) feed(c *gin.Context) {
	userId, err := getUserId(c, true)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var input PageInput
	c.BindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	output, err := h.services.IUsersPostsService.GetFeed(userId, input.Page)
	if err != nil {
		{
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

	}

	for ind, item := range output {
		file, err := os.OpenFile("user_data/profile_pictures/"+item.AuthorProfilePic, os.O_RDONLY, 0666)
		if err != nil {
			output[ind].AuthorProfilePic = c.Request.Host + "/images/base/base_pic.jpg"
		} else {
			output[ind].AuthorProfilePic = c.Request.Host + "/images/profiles/" + item.AuthorProfilePic
			file.Close()
		}
	}

	c.JSON(http.StatusOK, output)
}
