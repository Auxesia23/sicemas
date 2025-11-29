package services

import (
	"context"
	"situs-keagamaan/internal/app/repositories"
	"situs-keagamaan/internal/cache"
	"situs-keagamaan/internal/dto"
	"situs-keagamaan/internal/utils"
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
		return err
	}
	return nil
}
