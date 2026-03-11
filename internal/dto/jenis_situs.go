package dto

import "github.com/google/uuid"

type JenisSitusRequest struct {
	Nama      string `json:"nama" validate:"required"`
	Deskripsi string `json:"deskripsi" validate:"required"`
}

type JenisSitusResponse struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Nama        string    `json:"nama" db:"nama_jenis"`
	Deskripsi   string    `json:"deskripsi" db:"deskripsi"`
	JumlahSitus int       `json:"jumlah_situs" db:"jumlah_situs"`
}
