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

	CreateCase(ctx context.Context, cc models.CaseCreate) (int, error)
	GetCase(ctx context.Context, id int) (models.Case, error)
	GetCases(ctx context.Context, page, perPage int) ([]models.Case, error)
	UpdateCase(ctx context.Context, cu models.Case) error
	DeleteCase(ctx context.Context, id int) error
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
	CheckIfUserUsedRateById(ctx context.Context, rateID int) (bool, error)
	ClaimOutcome(ctx context.Context, userRateID, walletID int, amount float64) error
	ClaimDeposit(ctx context.Context, userRateID, walletID int, amount float64) error
	GetAll(ctx context.Context, userID, page, perPage int) ([]models.UserRateAdmin, error)
	UpdateNextDayCharge(ctx context.Context, userRateID int, nextDayCharge float64) error
}

type Wallet interface {
	Insert(ctx context.Context, userID int, token string, amount float64) error
	GetByUser(ctx context.Context, userID int) ([]models.Wallet, error)
	GetByToken(ctx context.Context, userID int, token string) (models.Wallet, error)
}

type User interface {
	Create(ctx context.Context, uc models.UserCreate) (int, error)
	GetByID(ctx context.Context, id int) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetHashedPwd(ctx context.Context, email string) (string, error)
	GetStatus(id int) (string, error)
	GetAll(ctx context.Context, page, perPage int, status string) ([]models.User, error)
	UpdateStatus(ctx context.Context, id int, status string) error
	UpdateProperties(ctx context.Context, id int, properties interface{}, startVerify bool) error
}

type Withdraw interface {
	Create(ctx context.Context, wu models.WithdrawCreate, walletID int) (int, error)
	GetByUserID(ctx context.Context, page, perPage, userID int) ([]models.Withdraw, error)
	GetByID(ctx context.Context, withdrawID int) (models.Withdraw, error)
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

type Voucher interface {
	GetVoucherByID(ctx context.Context, id string) (models.Voucher, error)
	// CreateVoucher properties must be a valid json
	CreateVoucher(ctx context.Context, id, voucherType string, properties []byte) error
	GetAllVouchers(ctx context.Context, offset, limit int) ([]models.Voucher, error)
	// UpdateVoucher properties must be a valid json
	UpdateVoucher(ctx context.Context, id, voucherType string, properties []byte) error
	DeleteVoucher(ctx context.Context, id string) error

	// Users vouchers
	CreateUserVoucher(ctx context.Context, userID int, voucherID string) error
	GetUserVoucher(ctx context.Context, userID int, voucherID string) (bool, error)
}
