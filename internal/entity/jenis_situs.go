package entity

import "github.com/google/uuid"

type JenisSitus struct {
	ID        uuid.UUID `db:"id"`
	NamaJenis string    `db:"nama_jenis"`
	Deskripsi string    `db:"deskripsi"`
}
