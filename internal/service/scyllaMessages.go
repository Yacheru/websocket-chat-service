package service

import (
	"context"
	"websocket-chat-service/internal/entities"
	"websocket-chat-service/internal/repository"
)

type ScyllaMessagesService struct {
	scylla repository.ScyllaRepository
}

func NewScyllaMessagesService(scylla repository.ScyllaRepository) *ScyllaMessagesService {
	return &ScyllaMessagesService{scylla: scylla}
}

func (s *ScyllaMessagesService) StoreMessageQuery(ctx context.Context, msg *entities.Message) error {
	return s.scylla.StoreMessageQuery(ctx, msg)
}
