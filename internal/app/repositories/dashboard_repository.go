package repositories

import (
	"context"
	"log"
	"sicemas/internal/dto"

	"github.com/jmoiron/sqlx"
)

type DashboardRepository interface {
	GetDashboardStats(ctx context.Context) (*dto.Stats, error)
	GetRecentSites(ctx context.Context) ([]dto.SitusKeagamaanResponse, error)
	GetStatistikJenis(ctx context.Context) ([]dto.StatistikJenis, error)
	GetRecentActivities(ctx context.Context) ([]dto.ActivityResponse, error)
}

type dashboardRepositoryImpl struct {
	DB *sqlx.DB
}

func NewDashboardRepo(db *sqlx.DB) DashboardRepository {
	return &dashboardRepositoryImpl{DB: db}
}

func (r *dashboardRepositoryImpl) GetDashboardStats(ctx context.Context) (*dto.Stats, error) {
	query := `
		SELECT
			(SELECT COUNT(id) FROM situs_keagamaan WHERE deleted_at IS NULL) AS total_situs,
			(SELECT COUNT(id) FROM jenis_situs) AS total_jenis,
			(SELECT COUNT(id) FROM users WHERE deleted_at IS NULL) AS petugas_aktif,
			(SELECT COUNT(id) FROM situs_keagamaan WHERE status_verifikasi = 'menunggu' AND deleted_at IS NULL) AS menunggu_verifikasi
	`

	var stats dto.Stats
	err := r.DB.GetContext(ctx, &stats, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &stats, nil
}

func (r *dashboardRepositoryImpl) GetRecentSites(ctx context.Context) ([]dto.SitusKeagamaanResponse, error) {
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
		LEFT JOIN users u ON s.pendata_id = u.id
		WHERE s.deleted_at IS NULL
		ORDER BY s.updated_at DESC
		LIMIT 5
	`

	var sites []dto.SitusKeagamaanResponse
	err := r.DB.SelectContext(ctx, &sites, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if sites == nil {
		return []dto.SitusKeagamaanResponse{}, nil
	}

	return sites, nil
}

func (r *dashboardRepositoryImpl) GetStatistikJenis(ctx context.Context) ([]dto.StatistikJenis, error) {
	query := `
		SELECT
			j.nama_jenis AS nama,
			COUNT(s.id) FILTER (WHERE s.deleted_at IS NULL) AS jumlah_situs
		FROM jenis_situs j
		LEFT JOIN situs_keagamaan s ON j.id = s.jenis_situs_id
		GROUP BY j.id, j.nama_jenis
		ORDER BY jumlah_situs DESC
	`

	var chartData []dto.StatistikJenis
	err := r.DB.SelectContext(ctx, &chartData, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if chartData == nil {
		return []dto.StatistikJenis{}, nil
	}

	return chartData, nil
}

func (r *dashboardRepositoryImpl) GetRecentActivities(ctx context.Context) ([]dto.ActivityResponse, error) {
	q := `
	SELECT
		a.id,
		u.nama_lengkap AS actor_name,
		a.action_type,
		a.target_name,
		a.entity_type,
		a.created_at
	FROM
		activities a
	JOIN
		users u ON a.user_id = u.id
	ORDER BY a.created_at DESC
	LIMIT 5;
	`
	var activities []dto.ActivityResponse
	if err := r.DB.SelectContext(ctx, &activities, q); err != nil {
		log.Println(err)
		return nil, err
	}
	return activities, nil
}
