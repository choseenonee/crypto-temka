package user

import (
	"context"
	"crypto-temka/internal/delivery/middleware/auth"
	"crypto-temka/internal/models"
	"crypto-temka/internal/repository"
	"crypto-temka/internal/service"
	"crypto-temka/pkg/log"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

const (
	HashCost = 10
)

type UserService struct {
	repo repository.UserRepository
	jwt  auth.JWTUtil
	log  *log.Logs
}

func InitUserService(
	repo repository.UserRepository,
	log *log.Logs,
	jwt auth.JWTUtil) service.UserService {
	return UserService{repo: repo, log: log, jwt: jwt}
}

func (serv UserService) CreateUser(ctx context.Context, user models.UserCreate) (int, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), HashCost)
	if err != nil {
		serv.log.Error(err.Error())
		return 0, "", err
	}

	newUser := models.UserCreate{
		UserBase: user.UserBase,
		Password: string(hash),
	}
	id, err := serv.repo.CreateUser(ctx, newUser)
	userToken := serv.jwt.CreateToken(id)
	if err != nil {
		serv.log.Error(err.Error())
		return 0, "", err
	}
	serv.log.Info(fmt.Sprintf("create user %v", id))
	return id, userToken, nil
}

func (serv UserService) Login(ctx context.Context, user models.UserLogin) (int, string, error) {
	id, pwd, err := serv.repo.GetHashPWD(ctx, user.Email)
	if err != nil {
		serv.log.Error(err.Error())
		return 0, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(user.Password))
	if err != nil {
		serv.log.Error(fmt.Sprintf("Invalide password to %v", id))
		return 0, "", fmt.Errorf("invalide password")
	}
	jwtToken := serv.jwt.CreateToken(id)
	serv.log.Info(fmt.Sprintf("login user %v", id))
	return id, jwtToken, nil
}

func (serv UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user, err := serv.repo.GetUserByID(ctx, id)
	if err != nil {
		serv.log.Error(err.Error())
		return nil, err
	}
	serv.log.Info(fmt.Sprintf("get user %v", id))
	return user, nil
}

func (serv UserService) Delete(ctx context.Context, id int) error {

	err := serv.repo.DeleteUser(ctx, id)
	if err != nil {
		serv.log.Error(err.Error())
		return err
	}
	serv.log.Info(fmt.Sprintf("delete user %v", id))
	return nil
}

func (serv UserService) AddPhoto(ctx context.Context, photo models.UserPhoto) error {
	err := serv.repo.AddPhoto(ctx, photo)
	if err != nil {
		serv.log.Error(err.Error())
		return err
	}
	serv.log.Info(fmt.Sprintf("add photo user %v", photo.ID))
	return nil
}

func (serv UserService) GetPhoto(ctx context.Context, idUser int) (*models.UserPhoto, error) {
	user, err := serv.repo.GetPhoto(ctx, idUser)
	if err != nil {
		serv.log.Error(err.Error())
		return nil, err
	}
	serv.log.Info(fmt.Sprintf("get photo user %v", idUser))
	return user, nil
}

func (serv UserService) SetStatus(ctx context.Context, status models.SetStatus) error {
	err := serv.repo.SetStatus(ctx, status)
	if err != nil {
		serv.log.Error(err.Error())
		return err
	}
	serv.log.Info(fmt.Sprintf("set status user %v", status.ID))
	return nil
}

func (serv UserService) Refresh(ctx context.Context, id int) (string, error) {
	userToken := serv.jwt.CreateToken(id)

	return userToken, nil
}
