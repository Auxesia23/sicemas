package services

import (
	"context"
	"fmt"
	"situs-keagamaan/internal/app/repositories"
	"situs-keagamaan/internal/cache"
	"situs-keagamaan/internal/dto"
	"situs-keagamaan/internal/utils"
	"time"
)

type AuthService interface {
	Login(ctx context.Context, in *dto.UserLogin) error
	VerifyOTP(ctx context.Context, in *dto.UserVerifyOTP) (*dto.Token, error)
}

type authServiceImpl struct {
	userRepo repositories.UserRepo
	cache    cache.Cache
}

func NewAuthService(userRepo repositories.UserRepo, cache cache.Cache) AuthService {
	return &authServiceImpl{
		userRepo,
		cache,
	}
}

func (s *authServiceImpl) Login(ctx context.Context, in *dto.UserLogin) error {
	index := utils.HashIndex(in.NIP)
	user, err := s.userRepo.ReadOne(ctx, index)
	if err != nil {
		return err
	}

	otp := utils.GenerateOTP6()
	s.cache.Set(ctx, fmt.Sprintf("otp:%v", user.ID), otp, 5*time.Minute)
	fmt.Printf("WA : %v | OTP : %v\nEMAIL : %v | OTP : %v\n", utils.Decrypt(user.NomorTelepon), otp, utils.Decrypt(user.Email), otp)

	return nil
}

func (s *authServiceImpl) VerifyOTP(ctx context.Context, in *dto.UserVerifyOTP) (*dto.Token, error) {
	index := utils.HashIndex(in.NIP)
	user, err := s.userRepo.ReadOne(ctx, index)
	if err != nil {
		return nil, err
	}
	var otp string
	if err := s.cache.Get(ctx, fmt.Sprintf("otp:%v", user.ID), &otp); err != nil {
		return nil, err
	}
	if otp != in.OTP {
		return nil, fmt.Errorf("OTP salah atau tidak ditemukan")
	}
	defer s.cache.Delete(ctx, fmt.Sprintf("otp:%v", user.ID))
	accesToken, _ := utils.GenerateAccessToken(user)
	token := &dto.Token{
		AccessToken:  accesToken,
		RefreshToken: accesToken,
	}
	return token, nil
}
