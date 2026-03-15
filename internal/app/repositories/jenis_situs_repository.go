package repositories

import (
	"context"
	"database/sql"
	"situs-keagamaan/internal/dto"
	"situs-keagamaan/internal/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type JenisSitusRepository interface {
	Create(ctx context.Context, in *dto.JenisSitusRequest) (*entity.JenisSitus, error)
	ReadAll(ctx context.Context) ([]dto.JenisSitusResponse, error)
	ReadOne(ctx context.Context, id uuid.UUID) (*dto.JenisSitusResponse, error)
	Update(ctx context.Context, id uuid.UUID, in *dto.JenisSitusRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type jenisSitusRepositoryImpl struct {
	DB *sqlx.DB
}

func NewJenisSitusRepo(db *sqlx.DB) JenisSitusRepository {
	return &jenisSitusRepositoryImpl{DB: db}
}

func (r *jenisSitusRepositoryImpl) Create(ctx context.Context, in *dto.JenisSitusRequest) (*entity.JenisSitus, error) {
	query := `INSERT INTO jenis_situs (nama_jenis, deskripsi) VALUES ($1, $2) RETURNING *`

	var newJenisSitus entity.JenisSitus
	if err := r.DB.GetContext(ctx, &newJenisSitus, query, in.Nama, in.Deskripsi); err != nil {
		return nil, err
	}
	return &newJenisSitus, nil
}

func (r *jenisSitusRepositoryImpl) ReadAll(ctx context.Context) ([]dto.JenisSitusResponse, error) {
	query := `
        SELECT
            js.id,
            js.nama_jenis,
            js.deskripsi,
            COUNT(s.id) AS jumlah_situs
        FROM jenis_situs js
        LEFT JOIN situs_keagamaan s ON js.id = s.jenis_situs_id
        GROUP BY js.id, js.nama_jenis, js.deskripsi
    `

	var jenisSitus []dto.JenisSitusResponse
	if err := r.DB.SelectContext(ctx, &jenisSitus, query); err != nil {
		return nil, err
	}
	return jenisSitus, nil
}

func (r *jenisSitusRepositoryImpl) ReadOne(ctx context.Context, id uuid.UUID) (*dto.JenisSitusResponse, error) {
	query := `
        SELECT
            js.id,
            js.nama_jenis,
            js.deskripsi,
            COUNT(s.id) AS jumlah_situs
        FROM jenis_situs js
        LEFT JOIN situs_keagamaan s ON js.id = s.jenis_situs_id
        WHERE js.id = $1
        GROUP BY js.id, js.nama_jenis, js.deskripsi
    `

	var jenisSitus dto.JenisSitusResponse
	if err := r.DB.GetContext(ctx, &jenisSitus, query, id); err != nil {
		return nil, err
	}
	return &jenisSitus, nil
}

func (r *jenisSitusRepositoryImpl) Update(ctx context.Context, id uuid.UUID, in *dto.JenisSitusRequest) error {
	query := `UPDATE jenis_situs SET nama_jenis = $1, deskripsi = $2 WHERE id = $3`

	result, err := r.DB.ExecContext(ctx, query, in.Nama, in.Deskripsi, id)
	if err != nil {
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

func (r *jenisSitusRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM jenis_situs WHERE id = $1`
	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
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
