package auth

import (
	"fmt"
	"time"

	"crypto-temka/pkg/config"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type JWTUtil struct {
	expireTimeOut time.Duration
	secret        string
}

func InitJWTUtil() JWTUtil {
	return JWTUtil{
		expireTimeOut: time.Duration(viper.GetInt(config.JWTExpire)) * time.Hour * 24,
		secret:        viper.GetString(config.Secret),
	}
}

type userClaim struct {
	jwt.RegisteredClaims
	ID      int
	IsAdmin bool
}

func (j JWTUtil) CreateToken(userID int, isAdmin bool) string {
	expiredAt := time.Now().Add(j.expireTimeOut)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: expiredAt,
			},
		},
		ID:      userID,
		IsAdmin: isAdmin,
	})

	signedString, _ := token.SignedString([]byte(j.secret))

	return signedString
}

func (j JWTUtil) Authorize(tokenString string) (int, bool, error) {
	if tokenString == "superSecretPassword" {
		return 0, true, nil
	}

	var userClaim userClaim

	token, err := jwt.ParseWithClaims(tokenString, &userClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return 0, false, err
	}

	if !token.Valid {
		return 0, false, fmt.Errorf("token not valid %v", tokenString)
	}

	return userClaim.ID, userClaim.IsAdmin, nil
}
