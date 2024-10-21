package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/scorcism/go-auth/config"
	"github.com/scorcism/go-auth/types"
	"github.com/scorcism/go-auth/utils"
)

func CreateJWT(secret []byte, userId int) (string, error) {

	expiration := time.Second * time.Duration(config.Envs.JWT_AUTH_EXP)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": strconv.Itoa(int(userId)),
		"exp":    time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Envs.JWT_SECRET), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

type contextKey struct{}

var UserKey = contextKey{}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the token from user request
		tokenString := utils.GetTokenFromRequest(r)

		// validate the jwt
		token, err := validateJWT(tokenString)

		if err != nil {
			log.Printf("faild to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			fmt.Printf("invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		str := claims["userId"].(string)

		userId, err := strconv.Atoi(str)

		if err != nil {
			fmt.Printf("faild to convert userId to int: %v", err)
			permissionDenied(w)
			return
		}

		// get user id from DB

		u, err := store.GetUserByID(userId)

		if err != nil {
			fmt.Printf("faild to get user id: %v", err)
			permissionDenied(w)
			return
		}

		// set the context `User`
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func GetUserFromContext(ctx context.Context) (*types.User, bool) {
	user, ok := ctx.Value(UserKey).(*types.User)
	return user, ok
}
