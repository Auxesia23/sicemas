package repositories

import (
	"context"
	"log"
	"situs-keagamaan/internal/dto"

	"github.com/jmoiron/sqlx"
)

type FotoSitusRepository interface {
	BulkCreate(ctx context.Context, data []dto.FotoSitusPayload) error
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
