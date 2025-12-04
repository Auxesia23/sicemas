package services

import (
	"context"
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/repositories"
	"situs-keagamaan/internal/cache"
	"situs-keagamaan/internal/dto"
	"situs-keagamaan/internal/utils"

	"github.com/lib/pq"
)

type UserService interface {
	Register(ctx context.Context, in *dto.UserRegister) error
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

func (s *userServiceImpl) Register(ctx context.Context, in *dto.UserRegister) error {
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
