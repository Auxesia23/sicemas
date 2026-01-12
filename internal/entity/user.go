package entity

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
