package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"websocket-chat-service/internal/server/websocket"
	"websocket-chat-service/internal/service"
)

type Handler struct {
	scylla service.ScyllaService
	ws     websocket.Client
}

func NewHandler(scylla service.ScyllaService, ws websocket.Client) *Handler {
	return &Handler{
		scylla: scylla,
		ws:     ws,
	}
}

func (h *Handler) Ws(ctx *gin.Context) {
	h.SuccessResponse(ctx, http.StatusOK, "hi", nil)
	return
}
