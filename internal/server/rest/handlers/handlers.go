package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"websocket-chat-service/internal/server/websocket"
	"websocket-chat-service/internal/service"
	"websocket-chat-service/pkg/constants"
)

type Handler struct {
	scylla    service.ScyllaService
	websocket websocket.Client
	wg        sync.WaitGroup
}

func NewHandler(scylla service.ScyllaService, websocket websocket.Client) *Handler {
	return &Handler{
		scylla:    scylla,
		websocket: websocket,
		wg:        sync.WaitGroup{},
	}
}

func (h *Handler) RunWS(ctx *gin.Context) {
	h.wg.Add(1)
	if err := h.websocket.Dial(ctx); err != nil {
		h.wg.Done()
		if errors.Is(err, constants.MaxLimitConnError) {
			h.ErrorResponse(ctx, http.StatusConflict, err.Error())
			return
		}

		h.ErrorResponse(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}
	h.wg.Done()

	h.SuccessResponse(ctx, http.StatusOK, "Success", nil)
	return
}
