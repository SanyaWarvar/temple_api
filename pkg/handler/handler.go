package handler

import (
	"net/http"

	"github.com/SanyaWarvar/temple_api/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	services *service.Service
	clients  map[string]*websocket.Conn
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services, clients: make(map[string]*websocket.Conn)}
}

func (h *Handler) InitRoutes(releaseMode bool) *gin.Engine {
	if releaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.HEAD("/health", h.check_health)

	router.GET("/ws/:access_token", h.ws)

	router.Static("/tik_toks", "./user_data/tik_toks")
	images := router.Group("/images")
	{
		images.Static("/profiles", "./user_data/profile_pictures")
		images.Static("/base", "./base_media")
	}

	auth := router.Group("/auth")
	{
		auth.POST("/sign_up", h.signUp)
		auth.POST("/sign_in", h.signIn)
		auth.POST("/send_code", h.sendConfirmCode)
		auth.POST("/confirm_email", h.confirmEmail)
		auth.POST("/refresh_token", h.refreshToken)
	}

	router.GET("/user/find", h.findUser)

	user := router.Group("/user", h.userIdentity)
	{
		user.GET("/:username", h.getUserInfo)
		user.PUT("/", h.updateUserInfo)
		user.PUT("/profile_pic", h.updateProfPic)   // new
		user.GET("/follows/:page", h.getAllFollows) // new
		user.GET("/subs/:page", h.getAllSubs)       // new
		friend := user.Group("/friends")
		{
			friend.GET("/:page", h.getAllFriends)
			friend.POST("/:username", h.inviteFriend)
			friend.DELETE("/:username", h.deleteFriend)
			friend.PUT("/:username", h.confirmFriend)
		}
		user.GET("/posts/:username", h.getPostsByU) // New

	}

	usersPosts := router.Group("/users_posts", h.userIdentity)
	{
		usersPosts.GET("/feed", h.feed)
		usersPosts.GET("/:id", h.getPost)
		usersPosts.POST("/", h.createPost)
		usersPosts.DELETE("/:id", h.deletePost)
		usersPosts.PUT("/:id", h.updatePost)

		usersPosts.PUT("/like/:id", h.likePost)
	}
	//new
	chats := router.Group("chats", h.userIdentity)
	{
		chats.GET("/", h.GetAllChats) //получить все чаты юзера
		chats.POST("/", h.CreateChat) //создать новый чат

		messages := chats.Group("/messages")
		{
			messages.GET("/:chat_id", h.GetChat)             // получить все сообщения из чата
			messages.POST("/:chat_id", h.NewMessage)         // отправить сообщение в чат
			messages.PUT("/read/:message_id", h.ReadMessage) // прочитать сообщение
			messages.PUT("/", h.EditMessage)                 // редактировать сообщение
			messages.DELETE("/:message_id", h.DeleteMessage) // удалить сообщение
		}

	}

	tiktoks := router.Group("rofls", h.userIdentity)
	{
		tiktoks.GET("/feed", h.feedTikTok)
		tiktoks.POST("/", h.createTiktok)
		tiktoks.GET("/:id", h.getTiktokById)
		tiktoks.DELETE("/:id", h.deleteTiktokById)

	}

	return router
}

func (h *Handler) check_health(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"details": "ok"})
}
