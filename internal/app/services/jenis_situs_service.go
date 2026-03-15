package services

import (
	"context"
	"database/sql"
	"errors"
	"log"
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/repositories"
	"situs-keagamaan/internal/dto"
	"situs-keagamaan/internal/entity"

	"github.com/google/uuid"
)

type JenisSitusService interface {
	CreateJenisSitus(ctx context.Context, in *dto.JenisSitusRequest, actorId uuid.UUID) (*dto.JenisSitusResponse, error)
	GetAllJenisSitus(ctx context.Context) ([]dto.JenisSitusResponse, error)
	UpdateJenisSitus(ctx context.Context, id uuid.UUID, in *dto.JenisSitusRequest, actorId uuid.UUID) error
	DeleteJenisSitus(ctx context.Context, id uuid.UUID, actorId uuid.UUID) error
}

type jenisSitusServiceImpl struct {
	repository   repositories.JenisSitusRepository
	activityRepo repositories.ActivityRepository
}

func NewJenisSitusService(repository repositories.JenisSitusRepository, activityRepo repositories.ActivityRepository) JenisSitusService {
	return &jenisSitusServiceImpl{repository: repository, activityRepo: activityRepo}
}

func (s *jenisSitusServiceImpl) CreateJenisSitus(ctx context.Context, in *dto.JenisSitusRequest, actorId uuid.UUID) (*dto.JenisSitusResponse, error) {
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
	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Menambahkan jenis-situs",
		EntityType: "JENIS-SITUS",
		EntityID:   newJenisSitus.ID.String(),
		TargetName: newJenisSitus.NamaJenis,
	})
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

func (s *jenisSitusServiceImpl) UpdateJenisSitus(ctx context.Context, id uuid.UUID, in *dto.JenisSitusRequest, actorId uuid.UUID) error {
	err := s.repository.Update(ctx, id, in)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFound("Jenis situs tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi kesalahan")
	}
	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Memperbarui data jenis-situs",
		EntityType: "JENIS-SITUS",
		EntityID:   id.String(),
		TargetName: in.Nama,
	})
	return nil
}

func (s *jenisSitusServiceImpl) DeleteJenisSitus(ctx context.Context, id, actorId uuid.UUID) error {
	targetJenisSitus, err := s.repository.ReadOne(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFound("Jenis situs tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi kesalahan")
	}

	err = s.repository.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFound("Jenis situs tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi kesalahan")
	}

	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Menghapus jenis-situs",
		EntityType: "JENIS-SITUS",
		EntityID:   id.String(),
		TargetName: targetJenisSitus.Nama,
	})

	return nil
}
