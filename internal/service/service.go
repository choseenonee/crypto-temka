package service

import (
	"context"
	"crypto-temka/internal/models"
)

type Static interface {
	CreateReview(ctx context.Context, rc models.ReviewCreate) (int, error)
	GetReviews(ctx context.Context, page, reviewsPerPage int) ([]models.Review, error)
	UpdateReview(ctx context.Context, r models.Review) error
	DeleteReview(ctx context.Context, id int) error

	SetMetrics(m models.MetricsSet) error
	GetMetrics() models.Metrics
}

type UserService interface {
	CreateUser(ctx context.Context, user models.UserCreate) (int, string, error)
	Login(ctx context.Context, user models.UserLogin) (int, string, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	AddPhoto(ctx context.Context, photo models.UserPhoto) error
	GetPhoto(ctx context.Context, idUser int) (*models.UserPhoto, error)
	SetStatus(ctx context.Context, status models.SetStatus) error
	Delete(ctx context.Context, id int) error
	Refresh(ctx context.Context, id int) (string, error)
}
