package repository

import (
	"context"
	"crypto-temka/internal/models"
)

type Static interface {
	CreateReview(ctx context.Context, rc models.ReviewCreate) (int, error)
	GetReviews(ctx context.Context, page, reviewsPerPage int) ([]models.Review, error)
	UpdateReview(ctx context.Context, r models.Review) error
	DeleteReview(ctx context.Context, id int) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user models.UserCreate) (int, error)
	GetHashPWD(ctx context.Context, email string) (int, string, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error

	AddPhoto(ctx context.Context, photo models.UserPhoto) error
	GetPhoto(ctx context.Context, idUser int) (*models.UserPhoto, error)

	SetStatus(ctx context.Context, status models.SetStatus) error
}
