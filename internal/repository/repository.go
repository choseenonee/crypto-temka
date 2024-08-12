package repository

import (
	"context"
	"crypto-temka/internal/models"
)

type Static interface {
	CreateReview(ctx context.Context, rc models.ReviewCreate) (int, error)
	GetReviews(ctx context.Context, page, perPage int) ([]models.Review, error)
	UpdateReview(ctx context.Context, r models.Review) error
	DeleteReview(ctx context.Context, id int) error
}

type Rate interface {
	CreateRate(ctx context.Context, rc models.RateCreate) (int, error)
	GetRates(ctx context.Context, page, perPage int) ([]models.Rate, error)
	UpdateRate(ctx context.Context, ru models.Rate) error
}
