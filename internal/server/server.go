package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"websocket-chat-service/internal/repository/redis"
	"websocket-chat-service/internal/server/websocket"

	"github.com/gin-gonic/gin"

	"websocket-chat-service/init/config"
	"websocket-chat-service/init/logger"
	"websocket-chat-service/internal/repository/scylladb"
	"websocket-chat-service/internal/server/rest/router"
	"websocket-chat-service/pkg/constants"
)

type HTTPServer struct {
	server *http.Server
	ws     websocket.Client
}

func NewServer(ctx context.Context, cfg *config.Config) (*HTTPServer, error) {
	scylla, err := scylladb.NewScyllaSession(ctx, cfg)
	if err != nil {
		return nil, err
	}

	_, err = redis.NewRedisClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	ws := websocket.NewWebSocket(cfg.WebsocketURL)

	engine := InitGinEngine(cfg)
	group := engine.Group(cfg.ApiEntry)
	router.NewRouterAndComponents(ctx, cfg, group, scylla).Router()

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.ApiPort),
		Handler:        engine,
		WriteTimeout:   10 * time.Second,
		ReadTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &HTTPServer{server: server, ws: ws}, nil
}

func (h *HTTPServer) Run(ctx context.Context) {
	go func() {
		if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(err.Error(), "ListenAndServe: "+constants.ServerLogger)
		}
	}()

	if err := h.ws.Dial(ctx); err != nil {
		logger.Error(err.Error(), "Dial: "+constants.ServerLogger)
		return
	}
}

func (h *HTTPServer) Shutdown(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}

func InitGinEngine(cfg *config.Config) *gin.Engine {
	var mode = gin.ReleaseMode
	if cfg.ApiDebug {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	engine := gin.New()

	engine.Use(gin.Recovery())
	engine.Use(gin.LoggerWithFormatter(logger.HTTPLogger))

	return engine
}
