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
	"github.com/lib/pq"
)

type RoleService interface {
	CreateRole(ctx context.Context, name string, actorId uuid.UUID) (*dto.RoleResponse, error)
	GetAllRole(ctx context.Context) ([]dto.RoleResponse, error)
	DeleteRole(ctx context.Context, id string, actorId uuid.UUID) error
}

type roleServiceImpl struct {
	roleRepo     repositories.RoleRepository
	activityRepo repositories.ActivityRepository
	logger       *slog.Logger
}

func NewRoleService(
	roleRepo repositories.RoleRepository,
	activityRepo repositories.ActivityRepository,
	logger *slog.Logger,
) RoleService {
	return &roleServiceImpl{
		roleRepo:     roleRepo,
		activityRepo: activityRepo,
		logger:       logger,
	}
}

func (s *roleServiceImpl) CreateRole(ctx context.Context, name string, actorId uuid.UUID) (*dto.RoleResponse, error) {
	s.logger.Info("creating new role", "name", name, "actor_id", actorId)

	newRole, err := s.roleRepo.Create(ctx, strings.ToLower(name))
	if err != nil {
		if e, ok := err.(*pq.Error); ok && e.Code == "23505" {
			s.logger.Warn("role already exists", "name", name)
			return nil, apperror.NewBadRequest("Role ini sudah ada.")
		}
		s.logger.Error("failed to create role", "error", err)
		return nil, apperror.NewInternal("Terjadi Kesalahan.")
	}

	response := dto.RoleResponse{
		ID:        newRole.ID,
		Name:      newRole.Name,
		CreatedAt: newRole.CreatedAt,
	}
	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Menambahkan role",
		EntityType: "ROLE",
		EntityID:   newRole.ID,
		TargetName: name,
	})

	s.logger.Info("role created successfully", "id", newRole.ID, "name", name)
	return &response, nil
}

func (s *roleServiceImpl) GetAllRole(ctx context.Context) ([]dto.RoleResponse, error) {
	s.logger.Info("fetching all roles")

	roles, err := s.roleRepo.ReadAll(ctx)
	if err != nil {
		s.logger.Error("failed to fetch roles", "error", err)
		return nil, apperror.NewInternal("Terjadi Kesalahan")
	}

	if len(roles) == 0 {
		s.logger.Info("no roles found")
		return []dto.RoleResponse{}, nil
	}

	var response []dto.RoleResponse
	for _, r := range roles {
		response = append(response, dto.RoleResponse{
			ID:        r.ID,
			Name:      r.Name,
			CreatedAt: r.CreatedAt,
		})
	}

	s.logger.Info("fetched roles successfully", "count", len(response))
	return response, nil
}

func (s *roleServiceImpl) DeleteRole(ctx context.Context, id string, actorId uuid.UUID) error {
	s.logger.Info("deleting role", "id", id, "actor_id", actorId)

	targetRole, err := s.roleRepo.ReadOne(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("role not found for deletion", "id", id)
			return apperror.NewNotFound("Role tidak ditemukan")
		}
		s.logger.Error("failed to read role for deletion", "error", err)
		return apperror.NewInternal("Terjadi Kesalahan")
	}

	if err := s.roleRepo.Delete(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			s.logger.Warn("role not found during deletion", "id", id)
			return apperror.NewNotFound("Role tidak ditemukan")
		}
		s.logger.Error("failed to delete role", "error", err)
		return apperror.NewInternal("Terjadi Kesalahan")
	}

	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Menghapus role",
		EntityType: "ROLE",
		EntityID:   targetRole.ID,
		TargetName: targetRole.Name,
	})

	s.logger.Info("role deleted successfully", "id", id)
	return nil
}
