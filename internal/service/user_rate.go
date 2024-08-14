package service

import (
	"context"
	"crypto-temka/internal/models"
	"crypto-temka/internal/repository"
	"crypto-temka/internal/utils"
	"crypto-temka/pkg/log"
	"errors"
	"fmt"
	"strings"
	"time"
)

type userRate struct {
	logger     *log.Logs
	repo       repository.UsersRate
	walletRepo repository.Wallet
	rateRepo   repository.Rate
}

func InitUserRate(repo repository.UsersRate, walletRepo repository.Wallet, rateRepo repository.Rate, logger *log.Logs) UserRate {
	return userRate{
		logger:     logger,
		repo:       repo,
		walletRepo: walletRepo,
		rateRepo:   rateRepo,
	}
}

func (u userRate) Create(ctx context.Context, urc models.UserRateCreate) (int, error) {
	rate, err := u.rateRepo.GetRate(ctx, urc.RateID)
	if err != nil {
		return 0, err
	}

	urc.Lock = utils.DateOnly(urc.Lock)
	rateLock := utils.DateOnly(time.Now()).Add(time.Hour * 24 * time.Duration(rate.MinLockDays))
	if urc.Lock.Before(rateLock) {
		return 0, fmt.Errorf("lock date must be more. min_lock_days for this rate is %v", rate.MinLockDays)
	}

	wallet, err := u.walletRepo.GetByToken(ctx, urc.UserID, urc.Token)
	if err != nil {
		if strings.HasPrefix(err.Error(), "sql: no rows in result set") {
			err := errors.New(fmt.Sprintf("wallet with token %v not found", urc.Token))
			u.logger.Error(err.Error())
			return 0, err
		}
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

func (u userRate) Get(ctx context.Context, id, userID int) (models.UserRate, error) {
	userRate, err := u.repo.Get(ctx, id)
	if err != nil {
		u.logger.Error(err.Error())
		return models.UserRate{}, err
	}

	if userRate.UserID != userID {
		return models.UserRate{}, fmt.Errorf("forbidden")
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

func (u userRate) Claim(ctx context.Context, userRateID, amount, userID int) error {
	userRate, err := u.repo.Get(ctx, userRateID)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	if userRate.UserID != userID {
		return fmt.Errorf("forbidden")
	}

	wallet, err := u.walletRepo.GetByToken(ctx, userRate.UserID, userRate.Token)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	err = u.repo.Claim(ctx, userRateID, amount, wallet.ID)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	return nil
}
