package scylladb

import (
	"context"
	"github.com/scylladb/gocqlx/v2"
	"websocket-chat-service/internal/entities"
)

type MessageRepository struct {
	scylla *gocqlx.Session
}

func NewMessageRepository(scylla *gocqlx.Session) *MessageRepository {
	return &MessageRepository{scylla}
}

func (m *MessageRepository) StoreMessage(ctx context.Context, msg *entities.Message) error {
	stmt := `INSERT INTO messages (player, player_uuid, message) VALUES (?, ?, ?)`

	if err := m.scylla.Session.Query(stmt, msg.Player.Username, msg.Player.UUID, msg.Message).WithContext(ctx).Exec(); err != nil {
		return err
	}

	return nil
}
