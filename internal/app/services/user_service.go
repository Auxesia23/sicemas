package services

import (
	"context"
	"database/sql"
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/repositories"
	"situs-keagamaan/internal/cache"
	"situs-keagamaan/internal/dto"
	"situs-keagamaan/internal/utils"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserService interface {
	Register(ctx context.Context, in *dto.UserRequest) error
	GetAllUser(ctx context.Context) ([]dto.UserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, in *dto.UserRequest) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type userServiceImpl struct {
	userRepo repositories.UserRepo
	cache    cache.Cache
}

func NewUserService(userRepo repositories.UserRepo, cache cache.Cache) UserService {
	return &userServiceImpl{
		userRepo,
		cache,
	}
}

func (s *userServiceImpl) Register(ctx context.Context, in *dto.UserRequest) error {
	index := utils.HashIndex(in.NIP)
	in.Email = utils.Encrypt(in.Email)
	in.NomorTelepon = utils.Encrypt(in.NomorTelepon)
	in.NIP = utils.Encrypt(in.NIP)

	if err := s.userRepo.Create(ctx, in, index); err != nil {
		e := err.(*pq.Error)
		if e.Code == "23505" {
			return apperror.NewBadRequest("User dengan NIP ini sudah ada.")
		}
		return apperror.NewInternal("Terjadi Kesalahan.")
	}
	return nil
}

func (s *userServiceImpl) GetAllUser(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := s.userRepo.ReadAll(ctx)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan")
	}

	var response []dto.UserResponse
	for _, u := range users {
		response = append(response, dto.UserResponse{
			ID:           u.ID,
			NIP:          utils.Decrypt(u.NIP),
			NamaLengkap:  u.NamaLengkap,
			Email:        utils.Decrypt(u.Email),
			NomorTelepon: utils.Decrypt(u.NomorTelepon),
			CreatedAt:    u.CreatedAt,
			UpdatedAt:    u.UpdatedAt,
		})
	}
	return response, nil
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, id uuid.UUID, in *dto.UserRequest) error {
	index := utils.HashIndex(in.NIP)
	in.Email = utils.Encrypt(in.Email)
	in.NomorTelepon = utils.Encrypt(in.NomorTelepon)
	in.NIP = utils.Encrypt(in.NIP)
	if err := s.userRepo.Update(ctx, id, in, index); err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFound("User tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi kesalahan")
	}
	return nil
}

func (s *userServiceImpl) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFound("User tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi kesalahan")
	}
	return nil
}
