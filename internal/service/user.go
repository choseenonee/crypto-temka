package service

import (
	"context"
	"crypto-temka/internal/delivery/middleware/auth"
	"crypto-temka/internal/models"
	"crypto-temka/internal/repository"
	"crypto-temka/pkg/log"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	logger *log.Logs
	repo   repository.User
	jwt    auth.JWTUtil
}

func InitUser(repo repository.User, jwt auth.JWTUtil, logger *log.Logs) User {
	return user{
		logger: logger,
		repo:   repo,
		jwt:    jwt,
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

	return user, nil
}

func (u user) GetAll(ctx context.Context, page, perPage int, status string) ([]models.User, error) {
	users, err := u.repo.GetAll(ctx, page, perPage, status)
	if err != nil {
		u.logger.Error(err.Error())
		return nil, err
	}

	return users, nil
}

func (u user) ChangeStatus(ctx context.Context, id int, status string) error {
	err := u.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	return nil
}
