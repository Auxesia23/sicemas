package services

import (
	"context"
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/repositories"
	"situs-keagamaan/internal/dto"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
)

type SitusKeagamaanService interface {
	CreateSitusKeagamaan(ctx context.Context, in *dto.SitusKeagamaanRequest, author uuid.UUID) error
	GetAllSitusKeagamaan(ctx context.Context) ([]dto.SitusKeagamaanResponse, error)
}

type situsKeagamaanServiceImpl struct {
	situsRepo repositories.SitusKeagamaanRepository
	enforcer  *casbin.Enforcer
}

func NewSitusKeagamaanService(situsRepo repositories.SitusKeagamaanRepository) SitusKeagamaanService {
	return &situsKeagamaanServiceImpl{
		situsRepo: situsRepo,
	}
}

func (s *situsKeagamaanServiceImpl) CreateSitusKeagamaan(ctx context.Context, in *dto.SitusKeagamaanRequest, author uuid.UUID) error {
	if len(in.Detail) == 0 {
		in.Detail = []byte("{}")
	}

	err := s.situsRepo.Create(ctx, in, author)
	if err != nil {
		return apperror.NewInternal("Terjadi kesalahan.")
	}
	return nil
}

func (s *situsKeagamaanServiceImpl) GetAllSitusKeagamaan(ctx context.Context) ([]dto.SitusKeagamaanResponse, error) {
	situs, err := s.situsRepo.ReadAll(ctx)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}
	if len(situs) == 0 {
		return []dto.SitusKeagamaanResponse{}, nil
	}
	return situs, nil
}
