package services

import (
	"context"
	"database/sql"
	"fmt"
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/repositories"
	"situs-keagamaan/internal/cache"
	"situs-keagamaan/internal/dto"
	"situs-keagamaan/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthService interface {
	Login(ctx context.Context, in *dto.UserLogin) error
	VerifyOTP(ctx context.Context, in *dto.UserVerifyOTP, loginContext *dto.SessionRequest) (*dto.Token, error)
	RefreshToken(ctx context.Context, refreshToken string, requestContext *dto.SessionRequest) (*dto.Token, error)
	Logout(ctx context.Context, refreshToken, accessToken string) error
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
		if err == sql.ErrNoRows {
			return apperror.NewNotFound("User dengan NIP ini tidak ditemukan.")
		}
		return apperror.NewInternal("Terjadi Kesalahan.")
	}

	otp := utils.GenerateOTP6()
	err = s.cache.Set(ctx, fmt.Sprintf("otp:%v", user.ID), otp, 5*time.Minute)
	if err != nil {
		return apperror.NewInternal("Terjadi kesalah.")
	}
	fmt.Println("OTP : ", otp)
	return nil
}

func (s *authServiceImpl) VerifyOTP(ctx context.Context, in *dto.UserVerifyOTP, loginContext *dto.SessionRequest) (*dto.Token, error) {
	index := utils.HashIndex(in.NIP)
	user, err := s.userRepo.ReadOne(ctx, index)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.NewNotFound("User dengan NIP ini tidak ditemukan.")
		}
		return nil, apperror.NewInternal("Terjadi Kesalahan.")
	}

	var otp string
	if err := s.cache.Get(ctx, fmt.Sprintf("otp:%v", user.ID), &otp); err != nil {
		if err == redis.Nil {
			return nil, apperror.NewUnauthorized("Tidak ada OTP untuk user ini.")
		}
		return nil, apperror.NewInternal("Terjadi Kesalahan.")
	}
	if otp != in.OTP {
		return nil, apperror.NewBadRequest("Kode OTP salah.")
	}

	refreshTokenId, err := uuid.NewV7()
	if err != nil {
		return nil, apperror.NewInternal("Terjadi Kesalahan.")
	}
	refreshToken, err := utils.GenerateRefreshToken(user, refreshTokenId)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi Kesalahan.")
	}

	accesTokenId, err := uuid.NewV7()
	if err != nil {
		return nil, apperror.NewInternal("Terjadi Kesalahan.")
	}
	accesToken, err := utils.GenerateAccessToken(user, refreshTokenId, accesTokenId)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi Kesalahan.")
	}

	sessionKey := fmt.Sprintf("rt:%v:%v", user.ID, refreshTokenId)
	sessionValue := &dto.SessionValue{
		UserID:      user.ID,
		SID:         refreshTokenId,
		UserAgent:   loginContext.UserAgent,
		IPAddress:   loginContext.IPAddress,
		GeoLocation: loginContext.GeoLocation,
		DeviceID:    loginContext.DeviceID,
	}
	err = s.cache.Set(ctx, sessionKey, sessionValue, 7*24*time.Hour)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi Kesalahan.")
	}

	token := &dto.Token{
		AccessToken:  accesToken,
		RefreshToken: refreshToken,
	}

	s.cache.Delete(ctx, fmt.Sprintf("otp:%v", user.ID))
	return token, nil
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, refreshToken string, requestContext *dto.SessionRequest) (*dto.Token, error) {
	claim, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, apperror.NewUnauthorized("Sesi tidak valid!")
	}

	user, err := s.userRepo.ReadById(ctx, claim.Subject)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.NewNotFound("Sesi tidak valid, user tidak ditemukan.")
		}
		return nil, apperror.NewInternal("Terjadi Kesalahan")
	}

	oldSessionKey := fmt.Sprintf("rt:%v:%v", claim.Subject, claim.ID)
	var oldSession dto.SessionValue
	err = s.cache.Get(ctx, oldSessionKey, &oldSession)
	if err != nil {
		if err == redis.Nil {
			return nil, apperror.NewUnauthorized("Sesi tidak valid!")
		}
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}

	trustScore := utils.CalculateTrustScore(requestContext, &oldSession)
	if trustScore <= 70 {
		return nil, apperror.NewUnauthorized("Terdeteksi perubahan signifikan pada sesi anda!")
	}

	newRefreshJti, err := uuid.NewV7()
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}
	newAccessJti, err := uuid.NewV7()
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}
	newSessionKey := fmt.Sprintf("rt:%v:%v", claim.Subject, newRefreshJti)
	newSession := &dto.SessionValue{
		UserID:      user.ID,
		SID:         newRefreshJti,
		UserAgent:   requestContext.UserAgent,
		IPAddress:   requestContext.IPAddress,
		GeoLocation: requestContext.GeoLocation,
		DeviceID:    requestContext.DeviceID,
	}

	err = s.cache.Set(ctx, newSessionKey, newSession, 7*24*time.Hour)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}

	newAccessToken, err := utils.GenerateAccessToken(user, newAccessJti, newRefreshJti)
	newRefreshToken, err := utils.GenerateRefreshToken(user, newRefreshJti)

	token := &dto.Token{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	err = s.cache.Delete(ctx, oldSessionKey)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi Kesalahan")
	}
	return token, nil
}

func (s *authServiceImpl) Logout(ctx context.Context, refreshToken, accessToken string) error {
	claim, _ := utils.ParseRefreshToken(refreshToken)
	if accessToken != "" {
		err := s.cache.Set(ctx, fmt.Sprintf("blocked:%v", accessToken), true, time.Minute*5)
		if err != nil {
			return apperror.NewInternal("Terjadi Kesalahan")
		}
	}
	err := s.cache.Delete(ctx, fmt.Sprintf("rt:%v:%v", claim.Subject, claim.ID))
	if err != nil {
		return apperror.NewInternal("Terjadi Kesalahan")
	}
	return nil
}
