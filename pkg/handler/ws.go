package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) ws(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer ws.Close()

	access_token := c.Param("access_token")
	accessToken, err := h.services.IJwtManagerService.ParseToken(access_token)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	userId := accessToken.UserId
	user, err := h.services.IUserService.GetUserById(userId)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	h.clients[user.Username] = ws
	fmt.Println(h.clients)

	ws.SetPongHandler(func(string) error {
		return ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	})

	ws.SetReadDeadline(time.Now().Add(60 * time.Second))

	go func() {
		for {
			time.Sleep(30 * time.Second)
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				logrus.Errorf("Error sending ping: %s", err)
				return
			}
		}
	}()

	message := "Welcome to the Temple Union!"
	if err := ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		logrus.Errorf("Error writing message: %s", err)
		return
	}

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			logrus.Errorf("Error reading message: %s", err)
			break
		}
	}

}
