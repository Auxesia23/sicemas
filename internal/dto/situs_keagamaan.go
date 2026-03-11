package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type SitusKeagamaanRequest struct {
	JenisSitusID  uuid.UUID `json:"jenis_situs_id" validate:"required"`
	KodeSitus     string    `json:"kode_situs"`
	Nama          string    `json:"nama" validate:"required"`
	JenisTipologi string    `json:"jenis_tipologi" validate:"required"`

	NomorTelepon    string `json:"nomor_telepon" validate:"required"`
	Email           string `json:"email"`
	Website         string `json:"website"`
	NomorBadanHukum string `json:"nomor_badan_hukum" validate:"required"`
	TahunBerdiri    int    `json:"tahun_berdiri" validate:"required"`

	Provinsi      string `json:"provinsi" validate:"required"`
	KabupatenKota string `json:"kabupaten_kota" validate:"required"`
	Kecamatan     string `json:"kecamatan" validate:"required"`
	Desa          string `json:"desa" validate:"required"`
	AlamatLengkap string `json:"alamat_lengkap" validate:"required"`

	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`

	LuasTanah            float64 `json:"luas_tanah" validate:"required"`
	LuasBangunan         float64 `json:"luas_bangunan" validate:"required"`
	StatusTanah          string  `json:"status_tanah" validate:"required"`
	NomorAIW             string  `json:"nomor_aiw" validate:"required"`
	NomorSertifikatWakaf string  `json:"nomor_sertifikat_wakaf" validate:"required"`
	DayaTampungMax       int     `json:"daya_tampung_max" validate:"required"`

	Detail json.RawMessage `json:"detail"`
}

type SitusKeagamaanResponse struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Nama           string    `json:"nama" db:"nama"`
	Jenis          string    `json:"jenis" db:"jenis"`
	Lokasi         string    `json:"lokasi" db:"lokasi"`
	Status         string    `json:"status" db:"status"`
	Pendata        string    `json:"pendata" db:"pendata"`
	TerakhirUpdate string    `json:"terakhir" db:"updated_at"`
}
