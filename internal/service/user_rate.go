package service

import (
	"context"
	"crypto-temka/internal/models"
	"crypto-temka/internal/repository"
	"crypto-temka/pkg/log"
)

type userRate struct {
	logger *log.Logs
	repo   repository.UsersRate
}

func (u userRate) Create(ctx context.Context, urc models.UserRateCreate) (int, error) {
	//TODO находим нужный воллетАйди, если его нет, то рейзим ошибку, также проверяем балик
	panic("implement me")
}

func (u userRate) Get(ctx context.Context, id int) (models.UserRate, error) {
	//TODO implement me
	panic("implement me")
}

func (u userRate) GetByUser(ctx context.Context, userID, page, perPage int) ([]models.UserRate, error) {
	//TODO implement me
	panic("implement me")
}

func InitUserRate(repo repository.UsersRate, logger *log.Logs) UserRate {
	return userRate{
		logger: logger,
		repo:   repo,
	}
}
