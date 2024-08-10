package service

import (
	"context"
	"crypto-temka/internal/models"
	"crypto-temka/internal/repository"
	"crypto-temka/pkg/log"
)

type static struct {
	logger *log.Logs
	repo   repository.Static
}

func InitStatic(repo repository.Static, logger *log.Logs) Static {
	return static{logger: logger, repo: repo}
}

func (s static) CreateReview(ctx context.Context, rc models.ReviewCreate) (int, error) {
	id, err := s.repo.CreateReview(ctx, rc)
	if err != nil {
		s.logger.Error(err.Error())
		return 0, err
	}

	return id, nil
}

func (s static) GetReviews(ctx context.Context, page, reviewsPerPage int) ([]models.Review, error) {
	reviews, err := s.repo.GetReviews(ctx, page, reviewsPerPage)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return reviews, nil
}

func (s static) UpdateReview(ctx context.Context, r models.Review) error {
	err := s.repo.UpdateReview(ctx, r)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	return nil
}

func (s static) DeleteReview(ctx context.Context, id int) error {
	err := s.repo.DeleteReview(ctx, id)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	return nil
}
