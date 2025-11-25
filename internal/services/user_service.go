package services

import (
	"context"
	"fmt"
	"situs-keagamaan/internal/models"
	"situs-keagamaan/internal/repositories"
	"situs-keagamaan/internal/utils"
	"time"
)

type UserService interface {
	RegisterUser(ctx context.Context, in *models.UserRegister) error
	LoginUser(ctx context.Context, in *models.UserLogin) error
	VerifyOTP(ctx context.Context, in *models.UserVerifyOTP) (*models.Token, error)
}

type userServiceImpl struct {
	userRepo repositories.UserRepo
	cache    repositories.Cache
}

func NewUserService(userRepo repositories.UserRepo, cache repositories.Cache) UserService {
	return &userServiceImpl{
		userRepo,
		cache,
	}
}

func (s *userServiceImpl) RegisterUser(ctx context.Context, in *models.UserRegister) error {
	index := utils.HashIndex(in.NIP)
	in.Email = utils.Encrypt(in.Email)
	in.NomorTelepon = utils.Encrypt(in.NomorTelepon)
	in.NIP = utils.Encrypt(in.NIP)

	if err := s.userRepo.Create(ctx, in, index); err != nil {
		return err
	}
	return nil
}

func (s *userServiceImpl) LoginUser(ctx context.Context, in *models.UserLogin) error {
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

func (s *userServiceImpl) VerifyOTP(ctx context.Context, in *models.UserVerifyOTP) (*models.Token, error) {
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
	token := &models.Token{
		AccessToken:  accesToken,
		RefreshToken: accesToken,
	}
	return token, nil
}
