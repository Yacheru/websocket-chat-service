package service

import "websocket-chat-service/internal/repository"

type ScyllaService interface {
}

type Service struct {
	ScyllaService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		ScyllaService: NewScyllaMessagesService(repo.ScyllaRepository),
	}
}
