package repositories

import (
	"context"
	"database/sql"
	"log"
	"sicemas/internal/dto"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type FotoSitusRepository interface {
	BulkCreate(ctx context.Context, data []dto.FotoSitusPayload) error
	GetPublicIDs(ctx context.Context, id []uuid.UUID) ([]string, error)
	GetBySitusID(ctx context.Context, situsId uuid.UUID) (*[]dto.FotoResponse, error)
	Delete(ctx context.Context, id []uuid.UUID, situsId uuid.UUID) error
}

type fotoSitusRepositoryImpl struct {
	DB *sqlx.DB
}

func NewFotoSitusRepo(db *sqlx.DB) FotoSitusRepository {
	return &fotoSitusRepositoryImpl{DB: db}
}

func (r *fotoSitusRepositoryImpl) BulkCreate(ctx context.Context, data []dto.FotoSitusPayload) error {
	if len(data) == 0 {
		return nil
	}

	query := `INSERT INTO foto_situs (situs_id, image_url, public_id)
              VALUES (:situs_id, :image_url, :public_id)`

	_, err := r.DB.NamedExecContext(ctx, query, data)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (r *fotoSitusRepositoryImpl) GetBySitusID(ctx context.Context, situsID uuid.UUID) (*[]dto.FotoResponse, error) {
	query := `SELECT id, image_url, public_id
              FROM foto_situs
              WHERE situs_id = $1`

	var result []dto.FotoResponse
	err := r.DB.SelectContext(ctx, &result, query, situsID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &result, nil
}

func (r *fotoSitusRepositoryImpl) GetPublicIDs(ctx context.Context, id []uuid.UUID) ([]string, error) {
	query := `SELECT public_id
              FROM foto_situs
              WHERE id = ANY($1)`

	var result []string
	err := r.DB.SelectContext(ctx, &result, query, pq.Array(id))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return result, nil
}

func (r *fotoSitusRepositoryImpl) Delete(ctx context.Context, id []uuid.UUID, situsId uuid.UUID) error {
	query := `DELETE FROM foto_situs
              WHERE id = ANY($1)
              AND situs_id = $2`

	result, err := r.DB.ExecContext(ctx, query, pq.Array(id), situsId)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
