package service

import (
	"context"
	"crypto-temka/internal/models"
	"crypto-temka/internal/repository"
	"crypto-temka/pkg/log"
)

type rate struct {
	logger *log.Logs
	repo   repository.Rate
}

func InitRate(repo repository.Rate, logger *log.Logs) Rate {
	return rate{
		logger: logger,
		repo:   repo,
	}
}

func (r rate) CreateRate(ctx context.Context, rc models.RateCreate) (int, error) {
	id, err := r.repo.CreateRate(ctx, rc)
	if err != nil {
		r.logger.Error(err.Error())
		return 0, err
	}

	return id, nil
}

func (r rate) GetRates(ctx context.Context, page, perPage int) ([]models.Rate, error) {
	rates, err := r.repo.GetRates(ctx, page, perPage)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, err
	}

	return rates, nil
}

func (r rate) UpdateRate(ctx context.Context, ru models.Rate) error {
	err := r.repo.UpdateRate(ctx, ru)
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}

	return nil
}
