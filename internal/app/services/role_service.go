package services

import (
	"context"
	"database/sql"
	"errors"
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
}

func NewRoleService(roleRepo repositories.RoleRepository, activityRepo repositories.ActivityRepository) RoleService {
	return &roleServiceImpl{
		roleRepo,
		activityRepo,
	}
}

func (s *roleServiceImpl) CreateRole(ctx context.Context, name string, actorId uuid.UUID) (*dto.RoleResponse, error) {
	newRole, err := s.roleRepo.Create(ctx, strings.ToLower(name))
	if err != nil {
		e := err.(*pq.Error)
		if e.Code == "23505" {
			return nil, apperror.NewBadRequest("Role ini sudah ada.")
		}
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
	return &response, nil
}

func (s *roleServiceImpl) GetAllRole(ctx context.Context) ([]dto.RoleResponse, error) {
	roles, err := s.roleRepo.ReadAll(ctx)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi Kesalahan")
	}

	if len(roles) == 0 {
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
	return response, nil
}

func (s *roleServiceImpl) DeleteRole(ctx context.Context, id string, actorId uuid.UUID) error {
	targetRole, err := s.roleRepo.ReadOne(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFound("Role tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi Kesalahan")
	}

	if err := s.roleRepo.Delete(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFound("Role tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi Kesalahan")
	}

	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Menghapus role",
		EntityType: "ROLE",
		EntityID:   targetRole.ID,
		TargetName: targetRole.Name,
	})
	return nil
}
