package dto

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Token struct {
	AccessToken  string
	RefreshToken string
	CSRFToken    string
}

type AccessToken struct {
	NamaLengkap string `json:"nama_lengkap"`
	SID         string `json:"SID"`
	jwt.RegisteredClaims
}

type SessionValue struct {
	UserID      uuid.UUID `json:"user_id"`
	SID         uuid.UUID `json:"sid"`
	UserAgent   string    `json:"user_agent"`
	IPAddress   string    `json:"ip_address"`
	GeoLocation string    `json:"geo_location"`
	DeviceID    string    `json:"device_id"`
}

type SessionRequest struct {
	UserAgent   string `json:"user_agent"`
	IPAddress   string `json:"ip_address"`
	GeoLocation string `json:"geo_location"`
	DeviceID    string `json:"device_id"`
}
