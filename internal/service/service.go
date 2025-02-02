package service

import (
	"context"
	"crypto-temka/internal/models"
)

type Static interface {
	// admin
	CreateReview(ctx context.Context, rc models.ReviewCreate) (int, error)
	UpdateReview(ctx context.Context, r models.Review) error
	DeleteReview(ctx context.Context, id int) error

	CreateCase(ctx context.Context, cc models.CaseCreate) (int, error)
	UpdateCase(ctx context.Context, cu models.Case) error
	DeleteCase(ctx context.Context, id int) error

	SetMetrics(m models.MetricsSet) error

	// public
	GetReviews(ctx context.Context, page, PerPage int) ([]models.Review, error)

	GetCase(ctx context.Context, id int) (models.Case, error)
	GetCases(ctx context.Context, page, perPage int) ([]models.Case, error)

	GetMetrics() models.Metrics

	GetOutcome() []models.Outcome
}

type Rate interface {
	// admin
	CreateRate(ctx context.Context, rc models.RateCreate) (int, error)
	UpdateRate(ctx context.Context, ru models.Rate) error

	// user
	GetRates(ctx context.Context, page, PerPage int) ([]models.Rate, error)
}

type UserRate interface {
	// user
	Create(ctx context.Context, urc models.UserRateCreate) (int, error)
	Get(ctx context.Context, id, userID int) (models.UserRate, error)
	GetByUser(ctx context.Context, userID, page, perPage int) ([]models.UserRate, error)
	ClaimOutcome(ctx context.Context, userRateID, userID int, amount float64) error
	ClaimDeposit(ctx context.Context, userRateID, userID int, amount float64) error

	// admin
	GetAll(ctx context.Context, userID, page, perPage int) ([]models.UserRateAdmin, error)
	UpdateNextDayCharge(ctx context.Context, userRateID int, nextDayCharge float64) error
}

type User interface {
	// public
	Register(ctx context.Context, uc models.UserCreate) (string, error)
	Auth(ctx context.Context, email, password string) (string, error)

	// user
	Get(ctx context.Context, id int) (models.User, error)
	GetStatus(id int) (string, error)
	UpdateProperties(ctx context.Context, id int, properties interface{}, startVerify bool) error

	// admin
	GetAll(ctx context.Context, page, perPage int, status string) ([]models.User, error)
	UpdateStatus(ctx context.Context, id int, status string) error
}

type Withdraw interface {
	// user
	Create(ctx context.Context, wu models.WithdrawCreate) (int, error)
	GetByUserID(ctx context.Context, page, perPage, userID int) ([]models.Withdraw, error)
	GetByID(ctx context.Context, withdrawID, userID int) (models.Withdraw, error)

	// admin
	GetAll(ctx context.Context, page, perPage int, status string) ([]models.Withdraw, error)
	UpdateStatus(ctx context.Context, withdrawID int, status string, properties interface{}) error
}

type Refer interface {
	// user
	Get(ctx context.Context, userID, page, perPage int) ([]models.Refer, error)
	Claim(ctx context.Context, id, userID int) error
}

type Message interface {
	Create(ctx context.Context, mc models.MessageCreate) (int, error)
	GetByID(ctx context.Context, id, userID int) (models.Message, error)
	GetByUser(ctx context.Context, userID int) ([]models.Message, error)
}
