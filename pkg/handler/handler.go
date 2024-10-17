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

func (h *Handler) InitRoutes(releaseMode bool) *gin.Engine {
	if releaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
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

	router.GET("/user/:username", h.getUserInfo)
	router.GET("/user/find", h.findUser)
	user := router.Group("/user", h.userIdentity)
	{
		user.PUT("/", h.updateUserInfo)
		friend := user.Group("/friends")
		{
			friend.GET("/", h.getAllFriends)
			friend.POST("/:username", h.inviteFriend)
			friend.DELETE("/:username", h.deleteFriend)
			friend.PUT("/:username", h.confirmFriend)
		}

	}

	return router
}

func (h *Handler) check_health(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"details": "ok"})
}
