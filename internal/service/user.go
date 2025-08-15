package service

import (
	"context"
	"fmt"

	"crypto-temka/internal/delivery/middleware/auth"
	"crypto-temka/internal/models"
	"crypto-temka/internal/repository"
	"crypto-temka/pkg/log"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	logger         *log.Logs
	repo           repository.User
	walletRepo     repository.Wallet
	withdrawalRepo repository.Withdraw
	userRateRepo   repository.UsersRate
	jwt            auth.JWTUtil
}

func InitUser(
	repo repository.User,
	walletRepo repository.Wallet,
	withdrawalRepo repository.Withdraw,
	userRateRepo repository.UsersRate,
	jwt auth.JWTUtil,
	logger *log.Logs,
) User {

	return user{
		logger:         logger,
		repo:           repo,
		walletRepo:     walletRepo,
		withdrawalRepo: withdrawalRepo,
		userRateRepo:   userRateRepo,
		jwt:            jwt,
	}
}

func (u user) Register(ctx context.Context, uc models.UserCreate) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(uc.Password), 10)
	if err != nil {
		u.logger.Error(err.Error())
		return "", err
	}

	uc.Password = string(hashedPwd)

	id, err := u.repo.Create(ctx, uc)
	if err != nil {
		u.logger.Error(err.Error())
		return "", err
	}

	token := u.jwt.CreateToken(id, false)

	return token, nil
}

func (u user) Auth(ctx context.Context, email, password string) (string, error) {
	hashedPwd, err := u.repo.GetHashedPwd(ctx, email)
	if err != nil {
		u.logger.Error(err.Error())
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(password)); err != nil {
		return "", fmt.Errorf("wrong password, err: %v", err.Error())
	}

	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		u.logger.Error(err.Error())
		return "", err
	}

	token := u.jwt.CreateToken(user.ID, false)

	return token, nil
}

func (u user) Get(ctx context.Context, id int) (models.User, error) {
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error(err.Error())
		return models.User{}, err
	}

	user.Wallets, err = u.walletRepo.GetByUser(ctx, id)
	if err != nil {
		u.logger.Error(err.Error())

		return models.User{}, err
	}

	user.Withdrawals, err = u.withdrawalRepo.GetByUserID(ctx, 1, 999999999999, user.ID)
	if err != nil {
		u.logger.Error(err.Error())

		return models.User{}, err
	}

	user.WalletsInsertHistory, err = u.walletRepo.GetWalletsInsertHistoryByUserID(ctx, user.ID)
	if err != nil {
		u.logger.Error(err.Error())

		return models.User{}, err
	}

	return user, nil
}

func (u user) GetAll(ctx context.Context, page, perPage int, status string) ([]models.User, error) {
	users, err := u.repo.GetAll(ctx, page, perPage, status)
	if err != nil {
		u.logger.Error(err.Error())
		return nil, err
	}

	for i := range users {
		users[i].Wallets, err = u.walletRepo.GetByUser(ctx, users[i].ID)
		if err != nil {
			u.logger.Error(err.Error())

			return nil, err
		}
	}

	for i := range users {
		users[i].Withdrawals, err = u.withdrawalRepo.GetByUserID(ctx, 1, 999999999999, users[i].ID)
		if err != nil {
			u.logger.Error(err.Error())

			return nil, err
		}
	}

	for i := range users {
		users[i].WalletsInsertHistory, err = u.walletRepo.GetWalletsInsertHistoryByUserID(ctx, users[i].ID)
		if err != nil {
			u.logger.Error(err.Error())

			return nil, err
		}
	}

	for i := range users {
		users[i].UserRates, err = u.userRateRepo.GetByUser(ctx, users[i].ID, 1, 99999999999)
		if err != nil {
			u.logger.Error(err.Error())

			return nil, err
		}
	}

	return users, nil
}

func (u user) GetStatus(id int) (string, error) {
	status, err := u.repo.GetStatus(id)
	if err != nil {
		u.logger.Error(err.Error())
		return "", err
	}

	return status, nil
}

func (u user) UpdateStatus(ctx context.Context, id int, status string) error {
	err := u.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	return nil
}

func (u user) UpdateProperties(ctx context.Context, id int, properties interface{}, startVerify bool) error {
	err := u.repo.UpdateProperties(ctx, id, properties, startVerify)
	if err != nil {
		u.logger.Error(err.Error())

		return err
	}

	return nil
}

func (u user) Update(ctx context.Context, user models.UserUpdate) error {
	err := u.repo.Update(ctx, user)
	if err != nil {
		u.logger.Error(err.Error())

		return err
	}

	return nil
}
