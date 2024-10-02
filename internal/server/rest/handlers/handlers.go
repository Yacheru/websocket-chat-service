package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"websocket-chat-service/internal/server/websocket"
	"websocket-chat-service/internal/service"
	"websocket-chat-service/pkg/constants"
)

type Handler struct {
	scylla    service.ScyllaService
	websocket websocket.Client
}

func NewHandler(scylla service.ScyllaService, websocket websocket.Client) *Handler {
	return &Handler{
		scylla:    scylla,
		websocket: websocket,
	}
}

func (h *Handler) RunWS(ctx *gin.Context) {
	if err := h.websocket.Dial(ctx); err != nil {
		if errors.Is(err, constants.MaxLimitConnError) {
			h.ErrorResponse(ctx, http.StatusConflict, err.Error())
			return
		}

		h.ErrorResponse(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.SuccessResponse(ctx, http.StatusOK, "Success", nil)
	return
}
