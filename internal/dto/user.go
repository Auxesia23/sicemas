package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserRequest struct {
	NIP          string  `json:"nip" validate:"required,min=16"`
	NamaLengkap  string  `json:"nama_lengkap" validate:"required"`
	Peran        string  `json:"peran" validate:"required"`
	Jabatan      string  `json:"jabatan" validate:"required"`
	UnitKerja    *string `json:"unit_kerja"`
	Email        string  `json:"email" validate:"omitempty,email"`
	NomorTelepon string  `json:"nomor_telepon" validate:"required"`
}

type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	NIP          string    `json:"nip"`
	NamaLengkap  string    `json:"nama_lengkap"`
	Jabatan      string    `json:"jabatan"`
	Peran        string    `json:"peran"`
	UnitKerja    string    `json:"unit_kerja"`
	Email        string    `json:"email"`
	NomorTelepon string    `json:"nomor_telepon"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserLogin struct {
	NIP string `json:"nip" validate:"required"`
}

type UserVerifyOTP struct {
	NIP string `json:"nip" validate:"required"`
	OTP string `json:"otp" validate:"required,min=6"`
}

type UserStepUpOTP struct {
	OTP string `json:"otp" validate:"required,min=6"`
}

type UserPermission struct {
	NamaLengkap string   `json:"nama_lengkap"`
	Permissions []string `json:"permissions"`
}
