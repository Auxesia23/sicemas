package repositories

import (
	"context"
	"log"
	"situs-keagamaan/internal/dto"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FotoSitusRepository interface {
	BulkCreate(ctx context.Context, data []dto.FotoSitusPayload) error
	GetBySitusID(ctx context.Context, situsID uuid.UUID) (*[]dto.FotoResponse, error)
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
