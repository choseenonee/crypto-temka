package service

import (
	"context"
	"crypto-temka/internal/models"
	"crypto-temka/internal/repository"
	"crypto-temka/pkg/log"
	"errors"
	"fmt"
)

type userRate struct {
	logger     *log.Logs
	repo       repository.UsersRate
	walletRepo repository.Wallet
}

func InitUserRate(repo repository.UsersRate, logger *log.Logs) UserRate {
	return userRate{
		logger: logger,
		repo:   repo,
	}
}

func (u userRate) Create(ctx context.Context, urc models.UserRateCreate) (int, error) {
	wallet, err := u.walletRepo.GetByToken(ctx, urc.UserID, urc.Token)
	if err != nil {
		u.logger.Error(err.Error())
		return 0, err
	}

	if wallet.ID == 0 {
		err := errors.New(fmt.Sprintf("wallet with token %v not found", urc.Token))
		u.logger.Error(err.Error())
		return 0, err
	}

	if wallet.Deposit < urc.Deposit {
		err = errors.New(fmt.Sprintf("not enough money on wallet with token %v", urc.Token))
		u.logger.Error(err.Error())
		return 0, err
	}

	id, err := u.repo.Create(ctx, urc, wallet.ID)
	if err != nil {
		u.logger.Error(err.Error())
		return 0, err
	}

	return id, nil
}

func (u userRate) Get(ctx context.Context, id int) (models.UserRate, error) {
	userRate, err := u.repo.Get(ctx, id)
	if err != nil {
		u.logger.Error(err.Error())
		return models.UserRate{}, err
	}

	return userRate, nil
}

func (u userRate) GetByUser(ctx context.Context, userID, page, perPage int) ([]models.UserRate, error) {
	userRates, err := u.repo.GetByUser(ctx, userID, page, perPage)
	if err != nil {
		u.logger.Error(err.Error())
		return nil, err
	}

	return userRates, nil
}

func (u userRate) ClaimOutcome(ctx context.Context, userRateID, amount int) error {
	userRate, err := u.repo.Get(ctx, userRateID)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	wallet, err := u.walletRepo.GetByToken(ctx, userRate.UserID, userRate.Token)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	err = u.repo.ClaimOutcome(ctx, userRateID, amount, wallet.ID)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	return nil
}

func (u userRate) ClaimDeposit(ctx context.Context, userRateID, amount int) error {
	userRate, err := u.repo.Get(ctx, userRateID)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	wallet, err := u.walletRepo.GetByToken(ctx, userRate.UserID, userRate.Token)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	err = u.repo.ClaimDeposit(ctx, userRateID, amount, wallet.ID)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	return nil
}
