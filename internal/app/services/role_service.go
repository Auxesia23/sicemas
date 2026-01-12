package services

import (
	"context"
	"database/sql"
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/repositories"
	"situs-keagamaan/internal/dto"

	"github.com/lib/pq"
)

type RoleService interface {
	CreateRole(ctx context.Context, name string) (*dto.RoleResponse, error)
	GetAllRole(ctx context.Context) ([]dto.RoleResponse, error)
	DeleteRole(ctx context.Context, id string) error
}

type roleServiceImpl struct {
	roleRepo repositories.RoleRepository
}

func NewRoleService(roleRepo repositories.RoleRepository) RoleService {
	return &roleServiceImpl{
		roleRepo,
	}
}

func (s *roleServiceImpl) CreateRole(ctx context.Context, name string) (*dto.RoleResponse, error) {
	newRole, err := s.roleRepo.Create(ctx, name)
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
	return &response, nil
}

func (s *roleServiceImpl) GetAllRole(ctx context.Context) ([]dto.RoleResponse, error) {
	roles, err := s.roleRepo.ReadAll(ctx)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi Kesalahan")
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

func (s *roleServiceImpl) DeleteRole(ctx context.Context, id string) error {
	if err := s.roleRepo.Delete(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFound("Role tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi Kesalahan")
	}
	return nil
}
