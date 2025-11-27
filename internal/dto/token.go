package dto

import (
	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	AccessToken  string
	RefreshToken string
}

type AccessToken struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}
