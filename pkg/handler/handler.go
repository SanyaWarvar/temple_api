package handler

import (
	"net/http"

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
	router.HEAD("/health", h.check_health)

	auth := router.Group("/auth")
	{
		auth.POST("/sign_up", h.signUp)
		auth.POST("/sign_in", h.signIn)
		auth.POST("/send_code", h.sendConfirmCode)
		auth.POST("/confirm_email", h.confirmEmail)
		auth.POST("/refresh_token", h.refreshToken)
	}

	user := router.Group("/user", h.userIdentity)
	{
		user.GET("/", h.getUserInfo)
		user.PUT("/", h.updateUserInfo)

	}

	return router
}

func (h *Handler) check_health(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"details": "ok"})
}
