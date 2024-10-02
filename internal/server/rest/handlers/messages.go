package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetAllMessages(ctx *gin.Context) {
	messages, err := h.scylla.GetAllMessages(ctx.Request.Context())
	if err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	h.SuccessResponse(ctx, http.StatusOK, "OK", messages)
}
