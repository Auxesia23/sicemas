package utils

import (
	"errors"
	"fmt"
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

func ParseAccessToken(tokenString string) (*jwt.RegisteredClaims, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("unexpected signing method")
			}
			return secret, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
