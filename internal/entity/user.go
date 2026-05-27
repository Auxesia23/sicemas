package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	NIPIndex     []byte    `db:"nip_index"`
	PasswordHash string    `db:"password_hash"`
	NamaLengkap  string    `db:"nama_lengkap"`
	Jabatan      string    `db:"jabatan"`
	UnitKerja    string    `db:"unit_kerja"`
	NIP          string    `db:"nip"`
	Email        string    `db:"email"`
	NomorTelepon string    `db:"nomor_telepon"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
