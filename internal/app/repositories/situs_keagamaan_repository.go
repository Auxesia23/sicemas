package repositories

import (
	"context"
	"log"
	"situs-keagamaan/internal/dto"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SitusKeagamaanRepository interface {
	Create(ctx context.Context, in *dto.SitusKeagamaanRequest, author uuid.UUID) error
	ReadAll(ctx context.Context) ([]dto.SitusKeagamaanResponse, error)
}

type situsKeagamaanRepositoryImpl struct {
	DB *sqlx.DB
}

func NewSitusKeagamaanRepo(db *sqlx.DB) SitusKeagamaanRepository {
	return &situsKeagamaanRepositoryImpl{DB: db}
}

func (r *situsKeagamaanRepositoryImpl) Create(ctx context.Context, in *dto.SitusKeagamaanRequest, author uuid.UUID) error {
	query := `
		INSERT INTO situs_keagamaan (
			jenis_situs_id, pendata_id, kode_situs, nama, jenis_tipologi,
			nomor_telepon, email, website, nomor_badan_hukum, tahun_berdiri,
			provinsi, kabupaten_kota, kecamatan, desa, alamat_lengkap, koordinat,
			luas_tanah, luas_bangunan, status_tanah, nomor_aiw, nomor_sertifikat_wakaf, daya_tampung_max,
			detail
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14,$15, ST_SetSRID(ST_MakePoint($16, $17), 4326),
			$18, $19, $20, $21, $22, $23,
			$24
		);`

	_, err := r.DB.ExecContext(ctx, query,
		in.JenisSitusID,         // $1
		author,                  // $2
		in.KodeSitus,            // $3
		in.Nama,                 // $4
		in.JenisTipologi,        // $5
		in.NomorTelepon,         // $6
		in.Email,                // $7
		in.Website,              // $8
		in.NomorBadanHukum,      // $9
		in.TahunBerdiri,         // $10
		in.Provinsi,             // $11
		in.KabupatenKota,        // $12
		in.Kecamatan,            // $13
		in.Desa,                 // $14
		in.AlamatLengkap,        // $15
		in.Longitude,            // $16 (PostGIS: Longitude duluan)
		in.Latitude,             // $17 (PostGIS: Latitude belakangan)
		in.LuasTanah,            // $18
		in.LuasBangunan,         // $19
		in.StatusTanah,          // $20
		in.NomorAIW,             // $21
		in.NomorSertifikatWakaf, // $22
		in.DayaTampungMax,       // $23
		in.Detail,               // $24
	)

	if err != nil {
		log.Println("Error insert situs keagamaan:", err.Error())
		return err
	}

	return nil
}

func (r *situsKeagamaanRepositoryImpl) ReadAll(ctx context.Context) ([]dto.SitusKeagamaanResponse, error) {
	query := `
		SELECT
			s.id,
			s.nama,
			j.nama_jenis AS jenis,
			s.desa AS lokasi,
			s.status_verifikasi AS status,
			u.nama_lengkap AS pendata,
			s.updated_at
		FROM situs_keagamaan s
		JOIN jenis_situs j ON s.jenis_situs_id = j.id
		JOIN users u ON s.pendata_id = u.id;
	`

	var situs []dto.SitusKeagamaanResponse
	err := r.DB.SelectContext(ctx, &situs, query)
	if err != nil {
		log.Println("Error read all situs keagamaan:", err.Error())
		return nil, err
	}

	return situs, nil
}
