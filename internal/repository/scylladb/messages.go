package scylladb

import (
	"context"
	"github.com/scylladb/gocqlx/v2"
	"time"
	"websocket-chat-service/init/logger"
	"websocket-chat-service/internal/entities"
	"websocket-chat-service/pkg/constants"
)

type MessageRepository struct {
	scylla *gocqlx.Session
}

func NewMessageRepository(scylla *gocqlx.Session) *MessageRepository {
	return &MessageRepository{scylla}
}

func (m *MessageRepository) GetAllMessages(ctx context.Context) ([]*entities.Message, error) {
	var entitiesMessages []*entities.Message

	stmt := `SELECT message, player, player_uuid, sent_at FROM messages`

	iter := m.scylla.ContextQuery(ctx, stmt, nil).Iter()

	var message string
	var playerUsername string
	var playerUUID string
	var sentAt = new(time.Time)

	for iter.Scan(&message, &playerUsername, &playerUUID, sentAt) {
		msg := &entities.Message{
			Message: message,
			Player: entities.Player{
				UUID:     playerUUID,
				Username: playerUsername,
			},
			SentAt: sentAt,
		}
		entitiesMessages = append(entitiesMessages, msg)
	}

	if err := iter.Close(); err != nil {
		logger.Error(err.Error(), constants.ScyllaCategory)
		return nil, err
	}

	return entitiesMessages, nil
}

func (m *MessageRepository) StoreMessage(ctx context.Context, msg *entities.Message) error {
	stmt := `INSERT INTO messages (player, player_uuid, message, sent_at) VALUES (?, ?, ?, ?)`

	if err := m.scylla.ContextQuery(ctx, stmt, nil).Bind(msg.Player.Username, msg.Player.UUID, msg.Message, msg.SentAt).Exec(); err != nil {
		logger.Error(err.Error(), constants.ScyllaCategory)
		return err
	}

	return nil
}
