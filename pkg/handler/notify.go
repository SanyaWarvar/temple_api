package handler

type NewMessageNotify1 struct {
	Type string            `json:"type"`
	Data NewMessageNotify2 `json:"data"`
}

type NewMessageNotify2 struct {
	AuthorUsername string `json:"from_username"`
	Message        string `json:"message"`
}

func (h *Handler) NewMessageNotify(authorUsername, recipientUsername, message string) {
	msg := NewMessageNotify1{
		Type: "new_message",
		Data: NewMessageNotify2{
			AuthorUsername: authorUsername,
			Message:        message,
		},
	}
	targetConn, ok := h.clients[recipientUsername]
	if !ok {
		return
	}

	targetConn.WriteJSON(msg)
}
