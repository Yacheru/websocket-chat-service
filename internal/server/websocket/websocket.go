package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/coder/websocket"

	"websocket-chat-service/init/logger"
	"websocket-chat-service/internal/entities"
	"websocket-chat-service/pkg/constants"
	"websocket-chat-service/pkg/utils"
)

type Client interface {
	Dial(ctx context.Context) error
}

type WebSocket struct {
	URL   string
	count int
	wg    *sync.WaitGroup
	mu    *sync.Mutex
	cChan chan *websocket.Conn
	mChan chan *entities.Message
}

func NewWebSocket(url string) *WebSocket {
	return &WebSocket{
		URL:   url,
		wg:    new(sync.WaitGroup),
		mu:    new(sync.Mutex),
		cChan: make(chan *websocket.Conn, 1),
		mChan: make(chan *entities.Message),
	}
}

func (ws *WebSocket) Dial(ctx context.Context) error {
	if ws.count == 1 {
		return errors.New("already connected")
	}

	ws.wg.Add(1)
	go func() {
		defer ws.wg.Done()

		c, _, err := websocket.Dial(ctx, ws.URL, nil)
		if err != nil {
			logger.Error(err.Error(), "Dial: "+constants.WebsocketLogger)
		}

		if err := c.Write(ctx, websocket.MessageText, []byte("Bearer auth")); err != nil {
			logger.Error(err.Error(), "Write auth: "+constants.WebsocketLogger)
			c.Close(websocket.StatusTryAgainLater, "try again later")
		}
		if err := c.Write(ctx, websocket.MessageText, []byte("Listen PlayerChatEvent")); err != nil {
			logger.Error(err.Error(), "Write event: "+constants.WebsocketLogger)
			c.Close(websocket.StatusTryAgainLater, "try again later")
		}

		ws.mu.Lock()
		ws.count++
		ws.mu.Unlock()

		ws.cChan <- c
	}()

	ws.wg.Wait()

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
						logger.Error("websocket connection closed", "Read: "+constants.WebsocketLogger)
						return
					}
					if errors.Is(err, io.EOF) {
						logger.Error("Connection closed by server", "Read: "+constants.WebsocketLogger)
						break
					}
					logger.Error(err.Error(), "WS Read: "+constants.WebsocketLogger)
					continue
				}

				if b[0] != 'E' {
					continue
				}

				bytes := utils.CutMessagePrefix(b)

				err = json.Unmarshal(bytes, message)
				if err != nil {
					logger.Error(err.Error(), constants.WebsocketLogger)
					continue
				}

				fmt.Println(message)
			}
		}
	}(<-ws.cChan)
}
