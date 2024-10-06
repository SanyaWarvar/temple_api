package handler

import (
	"github.com/SanyaWarvar/temple_api/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	auth := router.Group("/auth")
	{
		auth.POST("/sign_up", h.signUp)
		auth.POST("/sign_in", h.signIn)
		auth.POST("/send_code", h.sendConfirmCode)
		auth.POST("/confirm_email", h.confirmEmail)
		auth.POST("/refresh_token", h.refreshToken)
	}

	//api := router.Group("/api", h.userIdentity)

	return router
}
