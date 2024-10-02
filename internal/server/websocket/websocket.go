package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coder/websocket"
	"io"
	"sync"
	"websocket-chat-service/init/config"
	"websocket-chat-service/internal/service"

	"websocket-chat-service/init/logger"
	"websocket-chat-service/internal/entities"
	"websocket-chat-service/pkg/constants"
)

type Client interface {
	Dial(ctx context.Context) error
}

type WebSocket struct {
	url          string
	maxConnLimit int
	auth         string

	countConn int

	manager MessageManager

	wg *sync.WaitGroup
	mu sync.Mutex

	conn chan *websocket.Conn
}

func NewWebSocket(cfg *config.Config, service *service.Service) *WebSocket {
	manager := NewManager(service.ScyllaService)

	return &WebSocket{
		url:          cfg.WebsocketURL,
		maxConnLimit: cfg.WebsocketLimit,
		auth:         cfg.BearerAuth,
		manager:      manager,
		wg:           new(sync.WaitGroup),
		mu:           sync.Mutex{},
		conn:         make(chan *websocket.Conn, 1),
	}
}

func (ws *WebSocket) Dial(ctx context.Context) error {
	if ws.countConn == ws.maxConnLimit {
		return constants.MaxLimitConnError
	}

	ws.wg.Add(1)
	go func() {
		defer ws.wg.Done()

		c, _, err := websocket.Dial(ctx, ws.url, nil)
		if err != nil {
			logger.Error(err.Error(), "Dial: "+constants.WebsocketCategory)
		}

		if err := c.Write(ctx, websocket.MessageText, []byte(fmt.Sprintf("Bearer %s", ws.auth))); err != nil {
			logger.Error(err.Error(), "Write auth: "+constants.WebsocketCategory)
			c.Close(websocket.StatusTryAgainLater, "try again later")
		}
		if err := c.Write(ctx, websocket.MessageText, []byte("Listen PlayerChatEvent")); err != nil {
			logger.Error(err.Error(), "Write event: "+constants.WebsocketCategory)
			c.Close(websocket.StatusTryAgainLater, "try again later")
		}

		ws.mu.Lock()
		ws.countConn++
		ws.mu.Unlock()

		ws.conn <- c
	}()
	ws.wg.Wait()

	logger.Info("websocket connected", constants.WebsocketCategory)

	ws.listen(ctx)

	return nil
}

func (ws *WebSocket) listen(ctx context.Context) {
	go func(c *websocket.Conn) {
		defer func() {
			if err := c.Close(websocket.StatusNormalClosure, "connection closed"); err != nil {
			}
		}()

		var message = new(entities.Message)
		for {
			select {
			case <-ctx.Done():
				logger.Info(ctx.Err().Error(), "Close websocket connection")
				return
			default:
				_, b, err := c.Read(ctx)
				if err != nil {
					if websocket.CloseStatus(err) != -1 {
						logger.Error("websocket connection closed", "Read: "+constants.WebsocketCategory)
						return
					}
					if errors.Is(err, io.EOF) {
						logger.Error("Connection closed by server", "Read: "+constants.WebsocketCategory)
						break
					}
					logger.Error(err.Error(), "WS Read: "+constants.WebsocketCategory)
					continue
				}

				if b[0] != 'E' {
					continue
				}

				bytes := CutMessagePrefix(b)

				err = json.Unmarshal(bytes, message)
				if err != nil {
					logger.Error(err.Error(), constants.WebsocketCategory)
					continue
				}

				if err := ws.manager.ManageMessage(ctx, message); err != nil {
					continue
				}
			}
		}
	}(<-ws.conn)
}
