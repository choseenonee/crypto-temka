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
	GetRate(ctx context.Context, id int) (models.Rate, error)
}

type UsersRate interface {
	Create(ctx context.Context, urc models.UserRateCreate, walletID int) (int, error)
	Get(ctx context.Context, id int) (models.UserRate, error)
	GetByUser(ctx context.Context, userID, page, perPage int) ([]models.UserRate, error)
	ClaimOutcome(ctx context.Context, userRateID, amount, walletID int) error
	ClaimDeposit(ctx context.Context, userRateID, amount, walletID int) error
}

type Wallet interface {
	GetByUser(ctx context.Context, userID int) ([]models.Wallet, error)
	GetByToken(ctx context.Context, userID int, token string) (models.Wallet, error)
}
