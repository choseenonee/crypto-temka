package service

import (
	"context"
	"crypto-temka/internal/models"
	"crypto-temka/internal/repository"
	"crypto-temka/pkg/log"
)

type message struct {
	logger *log.Logs
	repo   repository.Message
}

func InitMessage(repo repository.Message, logger *log.Logs) Message {
	return message{
		logger: logger,
		repo:   repo,
	}
}

func (m message) Create(ctx context.Context, mc models.MessageCreate) (int, error) {
	id, err := m.repo.Create(ctx, mc)
	if err != nil {
		m.logger.Error(err.Error())
		return 0, err
	}

	return id, nil
}

func (m message) GetByID(ctx context.Context, id, userID int) (models.Message, error) {
	msg, err := m.repo.GetByID(ctx, id, userID)
	if err != nil {
		m.logger.Error(err.Error())
		return models.Message{}, err
	}

	return msg, nil
}

func (m message) GetByUser(ctx context.Context, userID int) ([]models.Message, error) {
	messages, err := m.repo.GetByUser(ctx, userID)
	if err != nil {
		m.logger.Error(err.Error())
		return nil, err
	}

	return messages, nil
}
