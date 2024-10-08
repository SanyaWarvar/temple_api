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

// TODO change 200 -> 201

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

	err := h.services.IUserService.CreateUser(input)
	if err != nil {
		errorMessage := ""
		if strings.Contains(err.Error(), "email") {
			errorMessage = "This email already exist"
		}
		if strings.Contains(err.Error(), "username") {
			errorMessage = "This username already exist"
		}
		newErrorResponse(c, http.StatusConflict, errorMessage)
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"details": "ok",
	})
}

func (h *Handler) sendConfirmCode(c *gin.Context) {
	var input models.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	target, err := h.services.IUserService.GetUserByUP(input)
	if err != nil || target.Email != input.Email {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	status, err := h.services.IEmailSmtpService.CheckEmailConfirm(input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if status {
		newErrorResponse(c, http.StatusConflict, "This email already confirmed")
		return
	}
	minTtl, _ := time.ParseDuration(os.Getenv("MIN_TTL"))
	maxTtl, _ := time.ParseDuration(os.Getenv("CODE_EXP"))

	_, ttl, err := h.services.ICacheService.GetConfirmCode(input.Email)
	if err == nil && minTtl < ttl {
		newErrorResponse(c, http.StatusConflict, fmt.Sprintf("Сode has already been sent %s ago", maxTtl-ttl))
		return
	}

	go h.services.IEmailSmtpService.SendConfirmEmailMessage(input)

	c.JSON(http.StatusCreated, map[string]interface{}{
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

	status, err := h.services.IEmailSmtpService.CheckEmailConfirm(input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "email not found")
		return
	}

	if status {
		newErrorResponse(c, http.StatusConflict, "This email already confirmed")
		return
	}

	err = h.services.IEmailSmtpService.ConfirmEmail(input.Email, input.Code)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

type SignInStruct struct {
	Email    string `json:"email" binding:"required" db:"email"`
	Password string `json:"password" binding:"required" db:"password_hash"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input SignInStruct

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	status, err := h.services.IEmailSmtpService.CheckEmailConfirm(input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "email not found")
		return
	}

	if !status {
		newErrorResponse(c, http.StatusForbidden, "This email not confirmed")
		return
	}

	target, err := h.services.IUserService.GetUserByEP(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Email != target.Email {
		newErrorResponse(c, http.StatusUnauthorized, "this user have another email")
		return
	}

	token, refresh, _, err := h.services.IJwtManagerService.GeneratePairToken(target.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]string{
		"access_token":  token,
		"refresh_token": refresh,
	})
}

type RefreshTokenInput struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *Handler) refreshToken(c *gin.Context) {
	var input RefreshTokenInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, err := h.services.IJwtManagerService.ParseToken(input.AccessToken)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "bad access token")
		return
	}
	expStatus := h.services.IJwtManagerService.CheckRefreshTokenExp(accessToken.RefreshId)
	if !expStatus {
		newErrorResponse(c, http.StatusUnauthorized, "refresh token is expired or not found")
		return
	}

	compareStatus := h.services.IJwtManagerService.CompareTokens(accessToken.RefreshId, input.RefreshToken)
	if !compareStatus {
		newErrorResponse(c, http.StatusBadRequest, "invalid refresh token")
		return
	}

	err = h.services.IJwtManagerService.DeleteRefreshTokenById(accessToken.RefreshId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, refresh, _, err := h.services.IJwtManagerService.GeneratePairToken(accessToken.UserId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]string{
		"access_token":  token,
		"refresh_token": refresh,
	})
}
