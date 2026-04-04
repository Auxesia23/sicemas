package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SitusKeagamaan struct {
	ID               uuid.UUID `db:"id"`
	StatusVerifikasi string    `db:"status_verifikasi"`
	JenisSitusID     uuid.UUID `db:"jenis_situs_id"`
	Pendata          uuid.UUID `db:"pendata"`

	KodeSitus     *string `db:"kode_situs"`
	Nama          *string `db:"nama"`
	JenisTipologi *string `db:"jenis_tipologi"`

	NomorTelepon    *string `db:"nomor_telepon"`
	Email           *string `db:"email"`
	Website         *string `db:"website"`
	NomorBadanHukum *string `db:"nomor_badan_hukum"`
	TahunBerdiri    *int    `db:"tahun_berdiri"`

	Provinsi      *string `db:"provinsi"`
	KabupatenKota *string `db:"kabupaten_kota"`
	Kecamatan     *string `db:"kecamatan"`
	Desa          *string `db:"desa"`
	AlamatLengkap *string `db:"alamat_lengkap"`

	Latitude  *float64 `db:"latitude" json:"latitude"`
	Longitude *float64 `db:"longitude" json:"longitude"`

	LuasTanah            *float64 `db:"luas_tanah"`
	LuasBangunan         *float64 `db:"luas_bangunan"`
	StatusTanah          *string  `db:"status_tanah"`
	NomorAIW             *string  `db:"nomor_aiw"`
	NomorSertifikatWakaf *string  `db:"nomor_sertifikat_wakaf"`
	DayaTampungMax       *int     `db:"daya_tampung_max"`

	Detail json.RawMessage `db:"detail"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
