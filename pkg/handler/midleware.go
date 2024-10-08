package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "Empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid auth header")
		return
	}

	accessToken, err := h.services.IJwtManagerService.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, accessToken.UserId)
}

func getUserId(c *gin.Context) (uuid.UUID, error) {
	userId, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return [16]byte{}, errors.New("user id not found")
	}

	idInt, ok := userId.(uuid.UUID)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id is invalid")
		return [16]byte{}, errors.New("user id is invalid")
	}

	return idInt, nil
}
