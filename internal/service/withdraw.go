package service

import (
	"context"
	"crypto-temka/internal/models"
	"crypto-temka/internal/repository"
	"crypto-temka/pkg/log"
	"fmt"
)

type withdraw struct {
	logger     *log.Logs
	repo       repository.Withdraw
	walletRepo repository.Wallet
}

func InitWithdraw(repo repository.Withdraw, walletRepo repository.Wallet, logger *log.Logs) Withdraw {
	return withdraw{
		logger:     logger,
		repo:       repo,
		walletRepo: walletRepo,
	}
}

// USER
func (w withdraw) Create(ctx context.Context, wu models.WithdrawCreate) (int, error) {
	walletID, err := w.walletRepo.GetByToken(ctx, wu.UserID, wu.Token)
	if err != nil {
		w.logger.Error(err.Error())
		return 0, err
	}

	id, err := w.repo.Create(ctx, wu, walletID.ID)
	if err != nil {
		w.logger.Error(err.Error())
		return 0, err
	}

	return id, nil
}

func (w withdraw) GetByUserID(ctx context.Context, page, perPage, userID int) ([]models.Withdraw, error) {
	withdrawals, err := w.repo.GetByUserID(ctx, page, perPage, userID)
	if err != nil {
		w.logger.Error(err.Error())
		return nil, err
	}

	return withdrawals, nil
}

func (w withdraw) GetByID(ctx context.Context, withdrawID, userID int) (models.Withdraw, error) {
	withdraw, err := w.repo.GetByID(ctx, withdrawID)
	if err != nil {
		w.logger.Error(err.Error())
		return models.Withdraw{}, err
	}

	if withdraw.UserID != userID {
		return models.Withdraw{}, fmt.Errorf("forbidden")
	}

	return withdraw, nil
}

// ADMIN
func (w withdraw) GetAll(ctx context.Context, page, perPage int, status string) ([]models.Withdraw, error) {
	withdrawals, err := w.repo.GetAll(ctx, page, perPage, status)
	if err != nil {
		w.logger.Error(err.Error())
		return nil, err
	}

	return withdrawals, nil
}

func (w withdraw) UpdateStatus(ctx context.Context, withdrawID int, status string, properties interface{}) error {
	err := w.repo.UpdateStatus(ctx, withdrawID, status, properties)
	if err != nil {
		w.logger.Error(err.Error())
		return err
	}

	return nil
}
