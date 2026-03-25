package repositories

import (
	"context"
	"database/sql"
	"log"
	"situs-keagamaan/internal/dto"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type SitusKeagamaanRepository interface {
	Create(ctx context.Context, in *dto.SitusKeagamaanRequest, author uuid.UUID) (uuid.UUID, error)
	ReadAll(ctx context.Context) ([]dto.SitusKeagamaanResponse, error)
	ReadOwn(ctx context.Context, userId uuid.UUID) ([]dto.SitusKeagamaanResponse, error)
	ReadDetail(ctx context.Context, id uuid.UUID) (*dto.SitusKeagamaanDetailResponse, error)
	ReadAllDetail(ctx context.Context, jenisSitus string) ([]dto.SitusKeagamaanDetailResponse, error)
	Update(ctx context.Context, id uuid.UUID, in *dto.SitusKeagamaanUpdate) error
	Delete(ctx context.Context, id uuid.UUID) error
	Verify(ctx context.Context, id uuid.UUID, in *dto.VerifikasiSitusRequest) error
	CheckOwnership(ctx context.Context, Id uuid.UUID, userId uuid.UUID) error
	ReadForPublic(ctx context.Context, filter dto.PublicListFilter) ([]dto.SitusPublicListResponse, error)
	ReadDetailForPublic(ctx context.Context, id uuid.UUID) (*dto.SitusPublicDetailResponse, error)
	GetLandingStats(ctx context.Context) (*dto.LandingStatsResponse, error)
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
			nomor_telepon, nomor_telpon_pengurus, email, website, nomor_badan_hukum, tahun_berdiri,
			provinsi, kabupaten_kota, kecamatan, desa, alamat_lengkap, koordinat,
			luas_tanah, luas_bangunan, status_tanah, nomor_aiw, nomor_sertifikat_wakaf, daya_tampung_max,
			detail
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, ST_SetSRID(ST_MakePoint($16, $17), 4326),
			$18, $19, $20, $21, $22, $23,
			$24
		) RETURNING id;`

	var id uuid.UUID
	err := r.DB.GetContext(ctx, &id, query,
		in.JenisSitusID,                  // $1
		author,                           // $2
		in.Nama,                          // $3
		in.JenisTipologi,                 // $4
		in.NomorTelepon,                  // $5
		pq.Array(in.NomorTelponPengurus), // $6
		in.Email,                         // $7
		in.Website,                       // $8
		in.NomorBadanHukum,               // $9
		in.TahunBerdiri,                  // $10
		in.Provinsi,                      // $11
		in.KabupatenKota,                 // $12
		in.Kecamatan,                     // $13
		in.Desa,                          // $14
		in.AlamatLengkap,                 // $15
		in.Longitude,                     // $16
		in.Latitude,                      // $17
		in.LuasTanah,                     // $18
		in.LuasBangunan,                  // $19
		in.StatusTanah,                   // $20
		in.NomorAIW,                      // $21
		in.NomorSertifikatWakaf,          // $22
		in.DayaTampungMax,                // $23
		in.Detail,                        // $24
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
				sk.nomor_telpon_pengurus,
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

	err := r.DB.GetContext(ctx, &detail, query, id)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &detail, nil
}

func (r *situsKeagamaanRepositoryImpl) ReadAllDetail(ctx context.Context, jenisSitus string) ([]dto.SitusKeagamaanDetailResponse, error) {
	query := `
			SELECT
				sk.id,
				js.nama_jenis AS jenis_situs,
				sk.status_verifikasi,
				sk.situs_id,
				sk.nama,
				sk.jenis_tipologi,
				sk.nomor_telepon,
				sk.nomor_telpon_pengurus,
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
			WHERE js.nama_jenis = $1
		`

	var detail []dto.SitusKeagamaanDetailResponse

	err := r.DB.SelectContext(ctx, &detail, query, jenisSitus)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return detail, nil
}

func (r *situsKeagamaanRepositoryImpl) Update(ctx context.Context, id uuid.UUID, in *dto.SitusKeagamaanUpdate) error {
	query := `
		UPDATE situs_keagamaan
		SET
			nama = $1,
			jenis_tipologi = $2,
			nomor_telepon = $3,
			nomor_telpon_pengurus = $4,
			email = $5,
			website = $6,
			nomor_badan_hukum = $7,
			tahun_berdiri = $8,
			provinsi = $9,
			kabupaten_kota = $10,
			kecamatan = $11,
			desa = $12,
			alamat_lengkap = $13,
			koordinat = ST_SetSRID(ST_MakePoint($14, $15), 4326),
			luas_tanah = $16,
			luas_bangunan = $17,
			status_tanah = $18,
			nomor_aiw = $19,
			nomor_sertifikat_wakaf = $20,
			daya_tampung_max = $21,
			detail = $22,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $23 AND status_verifikasi IN ('menunggu', 'ditolak')`

	result, err := r.DB.ExecContext(ctx, query,
		in.Nama,                          // $1
		in.JenisTipologi,                 // $2
		in.NomorTelepon,                  // $3
		pq.Array(in.NomorTelponPengurus), // $4
		in.Email,                         // $5
		in.Website,                       // $6
		in.NomorBadanHukum,               // $7
		in.TahunBerdiri,                  // $8
		in.Provinsi,                      // $9
		in.KabupatenKota,                 // $10
		in.Kecamatan,                     // $11
		in.Desa,                          // $12
		in.AlamatLengkap,                 // $13
		in.Longitude,                     // $14
		in.Latitude,                      // $15
		in.LuasTanah,                     // $16
		in.LuasBangunan,                  // $17
		in.StatusTanah,                   // $18
		in.NomorAIW,                      // $19
		in.NomorSertifikatWakaf,          // $20
		in.DayaTampungMax,                // $21
		in.Detail,                        // $22
		id,                               // $23
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

func (r *situsKeagamaanRepositoryImpl) ReadForPublic(ctx context.Context, filter dto.PublicListFilter) ([]dto.SitusPublicListResponse, error) {
	var query string
	var args []any

	if filter.UserLat != nil && filter.UserLng != nil {
		query = `
			SELECT
				s.id,
				s.nama,
				j.nama_jenis AS jenis,
				s.desa,

				ST_Y(s.koordinat::geometry) AS latitude,
				ST_X(s.koordinat::geometry) AS longitude,

				(SELECT image_url FROM foto_situs fs WHERE fs.situs_id = s.id ORDER BY created_at ASC LIMIT 1) AS foto_utama,

				-- Ngitung Jarak
				ST_DistanceSphere(
					s.koordinat,
					ST_SetSRID(ST_MakePoint($1, $2), 4326)
				) AS jarak_meter

			FROM situs_keagamaan s
			LEFT JOIN jenis_situs j ON s.jenis_situs_id = j.id
			WHERE s.status_verifikasi = 'terverifikasi'
			ORDER BY jarak_meter ASC;
		`
		args = append(args, *filter.UserLng, *filter.UserLat)

	} else {
		query = `
			SELECT
				s.id,
				s.nama,
				j.nama_jenis AS jenis,
				s.desa,

				ST_Y(s.koordinat::geometry) AS latitude,
				ST_X(s.koordinat::geometry) AS longitude,

				(SELECT image_url FROM foto_situs fs WHERE fs.situs_id = s.id ORDER BY created_at ASC LIMIT 1) AS foto_utama
				-- Gak ada kolom jarak_meter di sini

			FROM situs_keagamaan s
			LEFT JOIN jenis_situs j ON s.jenis_situs_id = j.id
			WHERE s.status_verifikasi = 'terverifikasi'
			ORDER BY s.nama ASC; -- Urut abjad aja
		`
	}

	var results []dto.SitusPublicListResponse

	err := r.DB.SelectContext(ctx, &results, query, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return results, nil
}

func (r *situsKeagamaanRepositoryImpl) ReadDetailForPublic(ctx context.Context, id uuid.UUID) (*dto.SitusPublicDetailResponse, error) {
	query := `
		SELECT
			sk.id,
			sk.nama,
			js.nama_jenis AS jenis,
			sk.alamat_lengkap,
			sk.desa,

			ST_Y(sk.koordinat::geometry) AS latitude,
			ST_X(sk.koordinat::geometry) AS longitude,

			COALESCE(sk.luas_tanah, 0) AS luas_tanah,
			COALESCE(sk.daya_tampung_max, 0) AS daya_tampung,

			-- FIX: Tambahin COALESCE biar kalau MT (ga ada fasilitas), DB balikin '{}' bukan NULL
			COALESCE(
				(
					SELECT jsonb_object_agg(key, value)
					FROM jsonb_each(sk.detail)
					WHERE key LIKE 'fasilitas_%'
				),
				'{}'::jsonb
			) AS fasilitas,

			COALESCE(
				(SELECT json_agg(fs.image_url) FROM foto_situs fs WHERE fs.situs_id = sk.id),
				'[]'::json
			) AS galeri

		FROM situs_keagamaan sk
		JOIN jenis_situs js ON sk.jenis_situs_id = js.id
		WHERE sk.id = $1 AND sk.status_verifikasi = 'terverifikasi';
	`

	var result dto.SitusPublicDetailResponse

	if err := r.DB.GetContext(ctx, &result, query, id); err != nil {
		log.Println("Error GetPublicDetail:", err)
		return nil, err
	}

	return &result, nil
}

func (r *situsKeagamaanRepositoryImpl) GetLandingStats(ctx context.Context) (*dto.LandingStatsResponse, error) {
	query := `
		SELECT
		    COUNT(id) AS total_situs,
		    COUNT(DISTINCT desa) AS total_desa,
		    COUNT(DISTINCT jenis_situs_id) AS total_kategori,
		    COALESCE(SUM(daya_tampung_max), 0) AS total_kapasitas
		FROM situs_keagamaan
		WHERE status_verifikasi = 'terverifikasi';
	`

	var result dto.LandingStatsResponse

	if err := r.DB.GetContext(ctx, &result, query); err != nil {
		log.Println("Error GetLandingStats:", err)
		return nil, err
	}

	return &result, nil
}
