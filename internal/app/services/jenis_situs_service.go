package services

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/repositories"
	"situs-keagamaan/internal/dto"
	"situs-keagamaan/internal/entity"
	"strings"

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
	logger       *slog.Logger
}

func NewJenisSitusService(
	repository repositories.JenisSitusRepository,
	activityRepo repositories.ActivityRepository,
	logger *slog.Logger,
) JenisSitusService {
	return &jenisSitusServiceImpl{
		repository:   repository,
		activityRepo: activityRepo,
		logger:       logger,
	}
}

func (s *jenisSitusServiceImpl) CreateJenisSitus(ctx context.Context, in *dto.JenisSitusRequest, actorId uuid.UUID) (*dto.JenisSitusResponse, error) {
	s.logger.Info("creating new jenis situs", "nama", in.Nama, "actor_id", actorId)

	newJenisSitus, err := s.repository.Create(ctx, in)
	if err != nil {
		s.logger.Error("failed to create jenis situs", "error", err)
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

	s.logger.Info("jenis situs created successfully", "id", newJenisSitus.ID, "nama", newJenisSitus.NamaJenis)
	return response, nil
}

func (s *jenisSitusServiceImpl) GetAllJenisSitus(ctx context.Context) ([]dto.JenisSitusResponse, error) {
	s.logger.Info("fetching all jenis situs")

	jenisSitus, err := s.repository.ReadAll(ctx)
	if err != nil {
		s.logger.Error("failed to fetch jenis situs", "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan")
	}
	if len(jenisSitus) == 0 {
		s.logger.Info("no jenis situs found")
		return []dto.JenisSitusResponse{}, nil
	}

	s.logger.Info("fetched jenis situs successfully", "count", len(jenisSitus))
	return jenisSitus, nil
}

func (s *jenisSitusServiceImpl) UpdateJenisSitus(ctx context.Context, id uuid.UUID, in *dto.JenisSitusRequest, actorId uuid.UUID) error {
	s.logger.Info("updating jenis situs", "id", id, "nama", in.Nama, "actor_id", actorId)

	err := s.repository.Update(ctx, id, in)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Warn("jenis situs not found for update", "id", id)
			return apperror.NewNotFound("Jenis situs tidak ditemukan")
		}
		s.logger.Error("failed to update jenis situs", "error", err)
		return apperror.NewInternal("Terjadi kesalahan")
	}

	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Memperbarui data jenis-situs",
		EntityType: "JENIS-SITUS",
		EntityID:   id.String(),
		TargetName: in.Nama,
	})

	s.logger.Info("jenis situs updated successfully", "id", id)
	return nil
}

func (s *jenisSitusServiceImpl) DeleteJenisSitus(ctx context.Context, id, actorId uuid.UUID) error {
	s.logger.Info("deleting jenis situs", "id", id, "actor_id", actorId)

	targetJenisSitus, err := s.repository.ReadOne(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("jenis situs not found for deletion", "id", id)
			return apperror.NewNotFound("Jenis situs tidak ditemukan")
		}
		s.logger.Error("failed to read jenis situs for deletion", "error", err)
		return apperror.NewInternal("Terjadi kesalahan")
	}

	err = s.repository.Delete(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "violates RESTRICT setting of foreign key constraint") || strings.Contains(err.Error(), "23503") {
			s.logger.Warn("jenis situs cannot be deleted, still used by other data", "id", id)
			return apperror.NewBadRequest("Gagal menghapus: Jenis situs ini tidak bisa dihapus karena masih digunakan oleh data situs keagamaan.")
		}
		if err == sql.ErrNoRows {
			s.logger.Warn("jenis situs not found during deletion", "id", id)
			return apperror.NewNotFound("Jenis situs tidak ditemukan")
		}
		s.logger.Error("failed to delete jenis situs", "error", err)
		return apperror.NewInternal("Terjadi kesalahan")
	}

	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Menghapus jenis-situs",
		EntityType: "JENIS-SITUS",
		EntityID:   id.String(),
		TargetName: targetJenisSitus.Nama,
	})

	s.logger.Info("jenis situs deleted successfully", "id", id)
	return nil
}
