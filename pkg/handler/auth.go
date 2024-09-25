package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input models.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid json")
		return
	}

	if valid := input.IsValid(); !valid {
		newErrorResponse(c, http.StatusBadRequest, "Invalid username or password")
		return
	}

	err := h.services.Authorizer.CreateUser(input)
	if err != nil {
		errorMessage := ""
		if strings.Contains(err.Error(), "email") {
			errorMessage = "This email already exist"
		}
		if strings.Contains(err.Error(), "username") {
			errorMessage = "This username already exist"
		}
		newErrorResponse(c, http.StatusBadRequest, errorMessage)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"details": "ok",
	})
}

func (h *Handler) sendConfirmCode(c *gin.Context) {
	var input models.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	status, err := h.services.Authorizer.CheckEmailConfirm(input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if status {
		newErrorResponse(c, http.StatusBadRequest, "This email already confirmed")
		return
	}
	minTtl, _ := time.ParseDuration(os.Getenv("MIN_TTL"))
	maxTtl, _ := time.ParseDuration(os.Getenv("CODE_EXP"))

	_, ttl, err := h.services.Authorizer.GetConfirmCode(input.Email)

	if err == nil && ttl > minTtl {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Сode has already been sent %s ago", ttl))
		return
	}

	go h.services.EmailSmtper.SendConfirmEmailMessage(input)

	c.JSON(http.StatusOK, map[string]interface{}{
		"exp_time":       maxTtl.String(),
		"next_code_time": (maxTtl - minTtl).String(),
	})
}

type ConfirmEmailInput struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

func (h *Handler) confirmEmail(c *gin.Context) {
	var input ConfirmEmailInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	status, err := h.services.Authorizer.CheckEmailConfirm(input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if status {
		newErrorResponse(c, http.StatusBadRequest, "This email already confirmed")
		return
	}

	err = h.services.Authorizer.ConfirmEmail(input.Email, input.Code)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"details": "ok",
	})
}
