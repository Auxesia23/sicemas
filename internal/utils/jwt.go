package utils

import (
	"errors"
	"os"
	"situs-keagamaan/internal/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(user *entity.User) (string, error) {
	var secret []byte
	secret = []byte(os.Getenv("JWT_SECRET"))

	claims := jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		Issuer:    "KUA Ciemas",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", errors.New("error signing token")
	}
	return tokenString, nil
}
