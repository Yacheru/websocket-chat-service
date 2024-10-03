package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetPlayerMessages(ctx *gin.Context) {
	nickname := ctx.Param("nickname")

	messages, err := h.scylla.GetPlayerMessages(ctx.Request.Context(), nickname)
	if err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if len(messages) == 0 {
		h.ErrorResponse(ctx, http.StatusNotFound, "Message not found")
		return
	}

	h.SuccessResponse(ctx, http.StatusOK, fmt.Sprintf("%s messages", nickname), messages)
	return
}

func (h *Handler) GetAllMessages(ctx *gin.Context) {
	messages, err := h.scylla.GetAllMessages(ctx.Request.Context())
	if err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	h.SuccessResponse(ctx, http.StatusOK, "OK", messages)
}
