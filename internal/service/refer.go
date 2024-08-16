package service

import (
	"context"
	"crypto-temka/internal/models"
	"crypto-temka/internal/repository"
	"crypto-temka/pkg/log"
)

type refer struct {
	logger *log.Logs
	repo   repository.Refer
}

func InitRefer(repo repository.Refer, logger *log.Logs) Refer {
	return refer{
		logger: logger,
		repo:   repo,
	}
}

func (r refer) Get(ctx context.Context, userID, page, perPage int) ([]models.Refer, error) {
	refers, err := r.repo.Get(ctx, userID, page, perPage)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, err
	}

	return refers, nil
}

func (r refer) Claim(ctx context.Context, id, userID int) error {
	err := r.repo.Claim(ctx, id, userID)
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}

	return nil
}
