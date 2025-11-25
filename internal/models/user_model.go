package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	NIPIndex     []byte    `db:"nip_index"`
	NamaLengkap  string    `db:"nama_lengkap"`
	NIP          string    `db:"nip"`
	Email        string    `db:"email"`
	NomorTelepon string    `db:"nomor_telepon"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type UserRegister struct {
	NIP          string `json:"nip" validate:"required,min=18"`
	NamaLengkap  string `json:"nama_lengkap" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	NomorTelepon string `json:"nomor_telepon" validate:"required"`
}

type UserResponse struct {
	ID           uuid.UUID `db:"id"`
	NIP          string    `db:"nip"`
	NamaLengkap  string    `db:"nama_lengkap"`
	Email        string    `db:"email"`
	NomorTelepon string    `db:"nomor_telepon"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type UserLogin struct {
	NIP string `json:"nip" validate:"required"`
}

type UserVerifyOTP struct {
	NIP string `json:"nip" validate:"required"`
	OTP string `json:"otp" validate:"required,min=6"`
}
