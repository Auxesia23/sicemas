package repositories

import (
	"context"
	"database/sql"
	"log"
	"situs-keagamaan/internal/dto"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SitusKeagamaanRepository interface {
	Create(ctx context.Context, in *dto.SitusKeagamaanRequest, author uuid.UUID) (uuid.UUID, error)
	ReadAll(ctx context.Context) ([]dto.SitusKeagamaanResponse, error)
	ReadOwn(ctx context.Context, userId uuid.UUID) ([]dto.SitusKeagamaanResponse, error)
	ReadDetail(ctx context.Context, id uuid.UUID) (*dto.SitusKeagamaanDetailResponse, error)
	Update(ctx context.Context, id uuid.UUID, in *dto.SitusKeagamaanUpdate) error
	Delete(ctx context.Context, id uuid.UUID) error
	Verify(ctx context.Context, id uuid.UUID, in *dto.VerifikasiSitusRequest) error
	CheckOwnership(ctx context.Context, Id uuid.UUID, userId uuid.UUID) error
}

type situsKeagamaanRepositoryImpl struct {
	DB *sqlx.DB
}

func NewSitusKeagamaanRepo(db *sqlx.DB) SitusKeagamaanRepository {
	return &situsKeagamaanRepositoryImpl{DB: db}
}

func (r *situsKeagamaanRepositoryImpl) Create(ctx context.Context, in *dto.SitusKeagamaanRequest, author uuid.UUID) (uuid.UUID, error) {
	query := `
        INSERT INTO situs_keagamaan (
            jenis_situs_id, pendata_id, nama, jenis_tipologi,
            nomor_telepon, email, website, nomor_badan_hukum, tahun_berdiri,
            provinsi, kabupaten_kota, kecamatan, desa, alamat_lengkap, koordinat,
            luas_tanah, luas_bangunan, status_tanah, nomor_aiw, nomor_sertifikat_wakaf, daya_tampung_max,
            detail
        ) VALUES (
            $1, $2, $3, $4,
            $5, $6, $7, $8, $9,
            $10, $11, $12, $13, $14, ST_SetSRID(ST_MakePoint($15, $16), 4326),
            $17, $18, $19, $20, $21, $22,
            $23
        ) RETURNING id;`

	var id uuid.UUID
	err := r.DB.GetContext(ctx, &id, query,
		in.JenisSitusID,         // $1
		author,                  // $2
		in.Nama,                 // $3
		in.JenisTipologi,        // $4
		in.NomorTelepon,         // $5
		in.Email,                // $6
		in.Website,              // $7
		in.NomorBadanHukum,      // $8
		in.TahunBerdiri,         // $9
		in.Provinsi,             // $10
		in.KabupatenKota,        // $11
		in.Kecamatan,            // $12
		in.Desa,                 // $13
		in.AlamatLengkap,        // $14
		in.Longitude,            // $15
		in.Latitude,             // $16
		in.LuasTanah,            // $17
		in.LuasBangunan,         // $18
		in.StatusTanah,          // $19
		in.NomorAIW,             // $20
		in.NomorSertifikatWakaf, // $21
		in.DayaTampungMax,       // $22
		in.Detail,               // $23
	)

	if err != nil {
		log.Println("Error insert situs keagamaan:", err.Error())
		return uuid.Nil, err
	}

	return id, nil
}

func (r *situsKeagamaanRepositoryImpl) ReadAll(ctx context.Context) ([]dto.SitusKeagamaanResponse, error) {
	query := `
        SELECT
            s.id,
            s.situs_id,
            s.nama,
            j.nama_jenis AS jenis,
            s.desa AS lokasi,
            s.status_verifikasi AS status,
            u.nama_lengkap AS pendata,
            s.updated_at
        FROM situs_keagamaan s
        LEFT JOIN jenis_situs j ON s.jenis_situs_id = j.id
        LEFT JOIN users u ON s.pendata_id = u.id;
    `

	var situs []dto.SitusKeagamaanResponse
	err := r.DB.SelectContext(ctx, &situs, query)
	if err != nil {
		log.Println("Error read all situs keagamaan:", err.Error())
		return nil, err
	}

	return situs, nil
}

func (r *situsKeagamaanRepositoryImpl) ReadOwn(ctx context.Context, userID uuid.UUID) ([]dto.SitusKeagamaanResponse, error) {
	query := `
		SELECT
			s.id,
			s.situs_id,
			s.nama,
			j.nama_jenis AS jenis,
			s.desa AS lokasi,
			s.status_verifikasi AS status,
			u.nama_lengkap AS pendata,
			s.updated_at
		FROM situs_keagamaan s
		JOIN jenis_situs j ON s.jenis_situs_id = j.id
		JOIN users u ON s.pendata_id = u.id
		WHERE s.pendata_id = $1;
	`

	var situs []dto.SitusKeagamaanResponse
	err := r.DB.SelectContext(ctx, &situs, query, userID)
	if err != nil {
		log.Println("Error read own situs keagamaan:", err.Error())
		return nil, err
	}

	return situs, nil
}

func (r *situsKeagamaanRepositoryImpl) ReadDetail(ctx context.Context, id uuid.UUID) (*dto.SitusKeagamaanDetailResponse, error) {
	query := `
			SELECT
				sk.id,
				js.nama_jenis AS jenis_situs,
				sk.status_verifikasi,
				sk.situs_id,
				sk.nama,
				sk.jenis_tipologi,
				sk.nomor_telepon,
				sk.email,
				sk.website,
				sk.nomor_badan_hukum,
				sk.tahun_berdiri,
				sk.provinsi,
				sk.kabupaten_kota,
				sk.kecamatan,
				sk.desa,
				sk.alamat_lengkap,
				ST_Y(sk.koordinat::geometry) AS latitude,
				ST_X(sk.koordinat::geometry) AS longitude,
				sk.luas_tanah,
				sk.luas_bangunan,
				sk.status_tanah,
				sk.nomor_aiw,
				sk.nomor_sertifikat_wakaf,
				sk.daya_tampung_max,
				sk.detail,
				sk.updated_at
			FROM situs_keagamaan sk
			LEFT JOIN jenis_situs js ON sk.jenis_situs_id = js.id
			WHERE sk.id = $1
		`

	var detail dto.SitusKeagamaanDetailResponse

	// Gunakan GetContext karena kita hanya mengambil 1 baris spesifik berdasarkan ID
	err := r.DB.GetContext(ctx, &detail, query, id)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &detail, nil
}

func (r *situsKeagamaanRepositoryImpl) Update(ctx context.Context, id uuid.UUID, in *dto.SitusKeagamaanUpdate) error {
	query := `
		UPDATE situs_keagamaan
		SET
			nama = $1,
			jenis_tipologi = $2,
			nomor_telepon = $3,
			email = $4,
			website = $5,
			nomor_badan_hukum = $6,
			tahun_berdiri = $7,
			provinsi = $8,
			kabupaten_kota = $9,
			kecamatan = $10,
			desa = $11,
			alamat_lengkap = $12,
			koordinat = ST_SetSRID(ST_MakePoint($13, $14), 4326),
			luas_tanah = $15,
			luas_bangunan = $16,
			status_tanah = $17,
			nomor_aiw = $18,
			nomor_sertifikat_wakaf = $19,
			daya_tampung_max = $20,
			detail = $21,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $22`

	result, err := r.DB.ExecContext(ctx, query,
		in.Nama,                 // $1
		in.JenisTipologi,        // $2
		in.NomorTelepon,         // $3
		in.Email,                // $4
		in.Website,              // $5
		in.NomorBadanHukum,      // $6
		in.TahunBerdiri,         // $7
		in.Provinsi,             // $8
		in.KabupatenKota,        // $9
		in.Kecamatan,            // $10
		in.Desa,                 // $11
		in.AlamatLengkap,        // $12
		in.Longitude,            // $13
		in.Latitude,             // $14
		in.LuasTanah,            // $15
		in.LuasBangunan,         // $16
		in.StatusTanah,          // $17
		in.NomorAIW,             // $18
		in.NomorSertifikatWakaf, // $19
		in.DayaTampungMax,       // $20
		in.Detail,               // $21
		id,                      // $22
	)

	if err != nil {
		log.Println("Error update situs keagamaan:", err.Error())
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
func (r *situsKeagamaanRepositoryImpl) CheckOwnership(ctx context.Context, situsID, userID uuid.UUID) error {
	query := `
		SELECT id
		FROM situs_keagamaan
		WHERE id = $1 AND pendata_id = $2
	`
	var result uuid.UUID
	err := r.DB.GetContext(ctx, &result, query, situsID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *situsKeagamaanRepositoryImpl) Verify(ctx context.Context, id uuid.UUID, in *dto.VerifikasiSitusRequest) error {
	query := `
		UPDATE situs_keagamaan
			SET
				status_verifikasi = $1,
				situs_id = NULLIF($2, ''),
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $3
	`
	_, err := r.DB.ExecContext(ctx, query, in.StatusVerifikasi, in.SitusID, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *situsKeagamaanRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM situs_keagamaan WHERE id = $1
	`
	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
