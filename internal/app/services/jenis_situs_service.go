package services

import (
	"context"
	"database/sql"
	"log"
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/repositories"
	"situs-keagamaan/internal/dto"

	"github.com/google/uuid"
)

type JenisSitusService interface {
	CreateJenisSitus(ctx context.Context, in *dto.JenisSitusRequest) (*dto.JenisSitusResponse, error)
	GetAllJenisSitus(ctx context.Context) ([]dto.JenisSitusResponse, error)
	UpdateJenisSitus(ctx context.Context, id uuid.UUID, in *dto.JenisSitusRequest) error
	DeleteJenisSitus(ctx context.Context, id uuid.UUID) error
}

type jenisSitusServiceImpl struct {
	repository repositories.JenisSitusRepository
}

func NewJenisSitusService(repository repositories.JenisSitusRepository) JenisSitusService {
	return &jenisSitusServiceImpl{repository: repository}
}

func (s *jenisSitusServiceImpl) CreateJenisSitus(ctx context.Context, in *dto.JenisSitusRequest) (*dto.JenisSitusResponse, error) {
	newJenisSitus, err := s.repository.Create(ctx, in)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan")
	}

	response := &dto.JenisSitusResponse{
		ID:          newJenisSitus.ID,
		Nama:        newJenisSitus.NamaJenis,
		Deskripsi:   newJenisSitus.Deskripsi,
		JumlahSitus: 0,
	}
	return response, nil
}

func (s *jenisSitusServiceImpl) GetAllJenisSitus(ctx context.Context) ([]dto.JenisSitusResponse, error) {
	jenisSitus, err := s.repository.ReadAll(ctx)
	if err != nil {
		log.Println(err.Error())
		return nil, apperror.NewInternal("Terjadi kesalahan")
	}
	if len(jenisSitus) == 0 {
		return []dto.JenisSitusResponse{}, nil
	}
	return jenisSitus, nil
}

func (s *jenisSitusServiceImpl) UpdateJenisSitus(ctx context.Context, id uuid.UUID, in *dto.JenisSitusRequest) error {
	err := s.repository.Update(ctx, id, in)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFound("Jenis situs tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi kesalahan")
	}
	return nil
}

func (s *jenisSitusServiceImpl) DeleteJenisSitus(ctx context.Context, id uuid.UUID) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFound("Jenis situs tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi kesalahan")
	}
	return nil
}
