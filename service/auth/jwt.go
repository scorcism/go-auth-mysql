package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/scorcism/go-auth/config"
)

func CreateJWT(secret []byte, userId int) (string, error) {
	expirationTime := time.Now().Add(time.Second * time.Duration(config.Envs.JWT_AUTH_EXP)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    expirationTime,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
