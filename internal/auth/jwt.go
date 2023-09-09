package auth

import (
	"time"
	"url-shortener/internal/config"
	"url-shortener/internal/model"

	"github.com/golang-jwt/jwt/v4"
)

func MakeJWT(user model.User) (string, error) {
	claims := model.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "atlantis",
			IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		},
		User: user,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Get().Auth.JWTSecretKey))
}
