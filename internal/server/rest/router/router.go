package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scylladb/gocqlx/v2"
	"websocket-chat-service/init/config"
	"websocket-chat-service/internal/repository"
	"websocket-chat-service/internal/server/rest/handlers"
	"websocket-chat-service/internal/server/websocket"
	"websocket-chat-service/internal/service"
)

type Router struct {
	router  *gin.RouterGroup
	handler *handlers.Handler
}

func NewRouterAndComponents(ctx context.Context, cfg *config.Config, router *gin.RouterGroup, scylla *gocqlx.Session) *Router {
	repo := repository.NewRepository(scylla)
	services := service.NewService(repo)
	ws := websocket.NewWebSocket(cfg.WebsocketURL)
	handler := handlers.NewHandler(services.ScyllaService, ws)

	return &Router{
		router:  router,
		handler: handler,
	}
}

func (r *Router) Router() {
	{
		r.router.GET("/ws", r.handler.Ws)
	}
}
