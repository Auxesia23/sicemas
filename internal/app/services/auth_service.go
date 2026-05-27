package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	apperror "sicemas/internal/app/appError"
	"sicemas/internal/app/repositories"
	"sicemas/internal/cache"
	"sicemas/internal/dto"
	"sicemas/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthService interface {
	Login(ctx context.Context, in *dto.UserLogin) error
	VerifyOTP(ctx context.Context, in *dto.UserVerifyOTP, loginContext *dto.SessionRequest) (*dto.Token, error)
	RefreshToken(ctx context.Context, refreshToken string, requestContext *dto.SessionRequest) (*dto.Token, error)
	Logout(ctx context.Context, refreshToken string, accessToken *string) error
	TriggerStepUpOTP(ctx context.Context, userId string) error
	VerifyStepUpOTP(ctx context.Context, userId string, sid uuid.UUID, userOtp string, loginContext *dto.SessionRequest) error
}

type authServiceImpl struct {
	userRepo repositories.UserRepo
	cache    cache.Cache
	logger   *slog.Logger
}

func NewAuthService(
	userRepo repositories.UserRepo,
	cache cache.Cache,
	logger *slog.Logger,
) AuthService {
	return &authServiceImpl{
		userRepo: userRepo,
		cache:    cache,
		logger:   logger,
	}
}

func (s *authServiceImpl) Login(ctx context.Context, in *dto.UserLogin) error {
	s.logger.Info("initiating login process")

	index := utils.HashIndex(in.NIP)
	user, err := s.userRepo.ReadOne(ctx, index)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Warn("user not found for login", "user_index", index)
			return apperror.NewNotFound("NIP atau password anda salah")
		}
		s.logger.Error("failed to read user from database", "user_index", index, "error", err)
		return apperror.NewInternal("Terjadi Kesalahan.")
	}

	if !utils.ComparePassword(user.PasswordHash, in.Password) {
		s.logger.Warn("invalid credentials for user", "user_id", user.ID)
		return apperror.NewUnauthorized("NIP atau password anda salah")
	}

	otp := utils.GenerateOTP6()
	err = s.cache.Set(ctx, fmt.Sprintf("otp:%v", user.ID), otp, 5*time.Minute)
	if err != nil {
		s.logger.Error("failed to store OTP in cache", "user_id", user.ID, "error", err)
		return apperror.NewInternal("Terjadi kesalahan.")
	}
	fmt.Println("OTP : ", otp)
	phoneNum, err := utils.Decrypt(user.NomorTelepon)
	if err != nil {
		s.logger.Error("failed to decrypt phone number", "user_id", user.ID, "error", err)
		return apperror.NewInternal("Terjadi kesalahan.")
	}
	if err := utils.SendWhatsAppOTP(phoneNum, otp); err != nil {
		s.logger.Error("failed to send WhatsApp OTP", "user_id", user.ID, "error", err)
		return apperror.NewInternal("Terjadi kesalahan.")
	}

	s.logger.Info("OTP sent successfully", "user_id", user.ID)
	return nil
}

func (s *authServiceImpl) VerifyOTP(ctx context.Context, in *dto.UserVerifyOTP, loginContext *dto.SessionRequest) (*dto.Token, error) {
	s.logger.Info("initiating OTP verification", "user_index", utils.HashIndex(in.NIP), "ip_address", loginContext.IPAddress)

	index := utils.HashIndex(in.NIP)
	user, err := s.userRepo.ReadOne(ctx, index)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Warn("user not found during OTP verification", "user_index", index)
			return nil, apperror.NewNotFound("User dengan NIP ini tidak ditemukan.")
		}
		s.logger.Error("failed to read user during OTP verification", "user_index", index, "error", err)
		return nil, apperror.NewInternal("Terjadi Kesalahan.")
	}

	var otp string
	if err := s.cache.Get(ctx, fmt.Sprintf("otp:%v", user.ID), &otp); err != nil {
		if err == redis.Nil {
			s.logger.Warn("no OTP found for user", "user_id", user.ID)
			return nil, apperror.NewUnauthorized("Tidak ada OTP untuk user ini.")
		}
		s.logger.Error("failed to get OTP from cache", "user_id", user.ID, "error", err)
		return nil, apperror.NewInternal("Terjadi Kesalahan.")
	}
	if otp != in.OTP {
		s.logger.Warn("incorrect OTP provided", "user_id", user.ID, "ip_address", loginContext.IPAddress)
		return nil, apperror.NewBadRequest("Kode OTP salah.")
	}

	refreshTokenId, err := uuid.NewV7()
	if err != nil {
		s.logger.Error("failed to generate refresh token ID", "user_id", user.ID, "error", err)
		return nil, apperror.NewInternal("Terjadi Kesalahan.")
	}
	accesTokenId, err := uuid.NewV7()
	if err != nil {
		s.logger.Error("failed to generate access token ID", "user_id", user.ID, "error", err)
		return nil, apperror.NewInternal("Terjadi Kesalahan.")
	}

	refreshToken, err := utils.GenerateRefreshToken(user, refreshTokenId)
	if err != nil {
		s.logger.Error("failed to generate refresh token", "user_id", user.ID, "error", err)
		return nil, apperror.NewInternal("Terjadi Kesalahan.")
	}
	accesToken, err := utils.GenerateAccessToken(user, accesTokenId, refreshTokenId)
	if err != nil {
		s.logger.Error("failed to generate access token", "user_id", user.ID, "error", err)
		return nil, apperror.NewInternal("Terjadi Kesalahan.")
	}
	csrfToken := utils.GenerateCSRFToken(user.ID)

	sessionKey := fmt.Sprintf("rt:%v:%v", user.ID, refreshTokenId)
	sessionValue := &dto.SessionValue{
		UserID:        user.ID,
		SID:           refreshTokenId,
		UserAgent:     loginContext.UserAgent,
		IPAddress:     loginContext.IPAddress,
		GeoLocation:   loginContext.GeoLocation,
		DeviceID:      loginContext.DeviceID,
		IsMFAVerified: true,
	}

	err = s.cache.Set(ctx, sessionKey, sessionValue, 7*24*time.Hour)
	if err != nil {
		s.logger.Error("failed to store session in cache", "user_id", user.ID, "session_id", refreshTokenId, "error", err)
		return nil, apperror.NewInternal("Terjadi Kesalahan saat menyimpan sesi.")
	}

	s.cache.Delete(ctx, fmt.Sprintf("otp:%v", user.ID))
	s.logger.Info("OTP verified successfully, session created", "user_id", user.ID, "jti", refreshTokenId.String(), "ip_address", loginContext.IPAddress)
	return &dto.Token{
		AccessToken:  accesToken,
		RefreshToken: refreshToken,
		CSRFToken:    csrfToken,
	}, nil
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, refreshToken string, requestContext *dto.SessionRequest) (*dto.Token, error) {
	s.logger.Info("initiating token refresh process", "ip_address", requestContext.IPAddress, "user_agent", requestContext.UserAgent)

	claim, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		s.logger.Warn("invalid or expired refresh token", "error", err, "ip_address", requestContext.IPAddress)
		return nil, apperror.NewUnauthorized("Sesi tidak valid atau telah kadaluarsa!")
	}

	user, err := s.userRepo.ReadById(ctx, claim.Subject)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Warn("user not found during token refresh", "user_id", claim.Subject)
			return nil, apperror.NewNotFound("Sesi tidak valid, user tidak ditemukan.")
		}
		s.logger.Error("failed to read user from database during token refresh", "user_id", claim.Subject, "error", err)
		return nil, apperror.NewInternal("Terjadi Kesalahan")
	}

	oldSessionKey := fmt.Sprintf("rt:%v:%v", claim.Subject, claim.ID)
	var oldSession dto.SessionValue
	err = s.cache.Get(ctx, oldSessionKey, &oldSession)
	if err != nil {
		if err == redis.Nil {
			s.logger.Warn("session not found or already expired", "user_id", claim.Subject, "session_id", claim.ID)
			return nil, apperror.NewUnauthorized("Sesi telah digunakan atau tidak ditemukan!")
		}
		s.logger.Error("failed to get session from cache during token refresh", "user_id", claim.Subject, "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}

	trustScore := utils.CalculateTrustScore(requestContext, &oldSession)
	if trustScore <= 70 {
		s.logger.Warn("suspicious activity detected, low trust score", "user_id", claim.Subject, "trust_score", trustScore, "ip_address", requestContext.IPAddress)
		_ = s.cache.Delete(ctx, oldSessionKey)
		return nil, apperror.NewUnauthorized("Terdeteksi aktivitas mencurigakan. Silakan login kembali.")
	}

	newRefreshJti, err := uuid.NewV7()
	if err != nil {
		s.logger.Error("failed to generate new refresh token ID", "user_id", claim.Subject, "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan sistem.")
	}
	newAccessJti, err := uuid.NewV7()
	if err != nil {
		s.logger.Error("failed to generate new access token ID", "user_id", claim.Subject, "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan sistem.")
	}

	newAccessToken, err := utils.GenerateAccessToken(user, newAccessJti, newRefreshJti)
	if err != nil {
		s.logger.Error("failed to generate new access token", "user_id", claim.Subject, "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}
	newRefreshToken, err := utils.GenerateRefreshToken(user, newRefreshJti)
	if err != nil {
		s.logger.Error("failed to generate new refresh token", "user_id", claim.Subject, "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}
	newCsrfToken := utils.GenerateCSRFToken(user.ID)

	err = s.cache.Delete(ctx, oldSessionKey)
	if err != nil {
		s.logger.Error("failed to delete old session from cache", "user_id", claim.Subject, "error", err)
		return nil, apperror.NewInternal("Terjadi Kesalahan saat membersihkan sesi lama.")
	}

	newSessionKey := fmt.Sprintf("rt:%v:%v", claim.Subject, newRefreshJti)
	newSession := &dto.SessionValue{
		UserID:        user.ID,
		SID:           newRefreshJti,
		UserAgent:     requestContext.UserAgent,
		IPAddress:     requestContext.IPAddress,
		GeoLocation:   requestContext.GeoLocation,
		DeviceID:      requestContext.DeviceID,
		IsMFAVerified: true,
	}

	err = s.cache.Set(ctx, newSessionKey, newSession, 7*24*time.Hour)
	if err != nil {
		s.logger.Error("failed to store new session in cache", "user_id", claim.Subject, "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan saat menyimpan sesi baru.")
	}

	s.logger.Info("token refreshed successfully", "user_id", claim.Subject, "new_jti", newRefreshJti.String())
	return &dto.Token{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		CSRFToken:    newCsrfToken,
	}, nil
}

func (s *authServiceImpl) Logout(ctx context.Context, refreshToken string, accessToken *string) error {
	s.logger.Info("initiating logout process")

	claim, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		s.logger.Warn("invalid refresh token during logout", "error", err)
		return apperror.NewBadRequest("Token tidak valid.")
	}

	if accessToken != nil && *accessToken != "" {
		err := s.cache.Set(ctx, fmt.Sprintf("blocked:%v", *accessToken), true, time.Minute*15)
		if err != nil {
			s.logger.Error("failed to block access token", "user_id", claim.Subject, "error", err)
			return apperror.NewInternal("Terjadi Kesalahan saat blokir akses token")
		}
		s.logger.Info("access token blocked for security", "user_id", claim.Subject, "jti", claim.ID)
	}

	sessionKey := fmt.Sprintf("rt:%v:%v", claim.Subject, claim.ID)
	err = s.cache.Delete(ctx, sessionKey)
	if err != nil {
		s.logger.Error("failed to delete session from cache during logout", "user_id", claim.Subject, "error", err)
		return apperror.NewInternal("Terjadi Kesalahan saat menghapus sesi Redis")
	}

	s.logger.Info("user logged out successfully, session terminated", "user_id", claim.Subject, "jti", claim.ID)
	return nil
}

func (s *authServiceImpl) TriggerStepUpOTP(ctx context.Context, userId string) error {
	s.logger.Info("initiating StepUp verification", "user_id", userId)
	user, err := s.userRepo.ReadById(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("user not found for verification", "user_id", userId)
			return apperror.NewNotFound("User tidak ditemukan")
		}
		s.logger.Error("failed to read user for verification", "error", err)
		return apperror.NewInternal("Terjadi kesalahan saat mencari data user")
	}

	otp := utils.GenerateOTP6()
	err = s.cache.Set(ctx, fmt.Sprintf("stepup_otp:%v", user.ID), otp, 5*time.Minute)
	if err != nil {
		s.logger.Error("failed to store OTP in cache", "user_id", user.ID, "error", err)
		return apperror.NewInternal("Terjadi kesalahan.")
	}
	fmt.Println("OTP : ", otp)
	phoneNum, err := utils.Decrypt(user.NomorTelepon)
	if err != nil {
		s.logger.Error("failed to decrypt phone number", "user_id", user.ID, "error", err)
		return apperror.NewInternal("Terjadi kesalahan.")
	}
	if err := utils.SendWhatsAppOTP(phoneNum, otp); err != nil {
		s.logger.Error("failed to send WhatsApp OTP", "user_id", user.ID, "error", err)
		return apperror.NewInternal("Terjadi kesalahan.")
	}
	return nil
}

func (s *authServiceImpl) VerifyStepUpOTP(ctx context.Context, userId string, sid uuid.UUID, userOtp string, loginContext *dto.SessionRequest) error {
	s.logger.Info("verifying StepUp OTP", "user_id", userId)

	otpKey := fmt.Sprintf("stepup_otp:%v", userId)
	var otp string

	if err := s.cache.Get(ctx, otpKey, &otp); err != nil {
		if err == redis.Nil {
			s.logger.Warn("no StepUp OTP found for user or expired", "user_id", userId)
			return apperror.NewUnauthorized("OTP tidak ditemukan atau sudah kadaluarsa. Silakan minta ulang.")
		}
		s.logger.Error("failed to get OTP from cache", "user_id", userId, "error", err)
		return apperror.NewInternal("Terjadi Kesalahan sistem.")
	}

	if otp != userOtp {
		s.logger.Warn("incorrect StepUp OTP provided", "user_id", userId, "ip_address", loginContext.IPAddress)
		return apperror.NewBadRequest("Kode OTP salah.")
	}

	sessionKey := fmt.Sprintf("rt:%v:%v", userId, sid)
	var currentSession dto.SessionValue

	if err := s.cache.Get(ctx, sessionKey, &currentSession); err != nil {
		s.logger.Error("failed to get session from cache during step up", "user_id", userId, "error", err)
		return apperror.NewUnauthorized("Sesi tidak valid atau sudah kadaluarsa.")
	}

	currentSession.IsMFAVerified = true
	currentSession.IPAddress = loginContext.IPAddress
	currentSession.UserAgent = loginContext.UserAgent
	currentSession.GeoLocation = loginContext.GeoLocation
	currentSession.DeviceID = loginContext.DeviceID

	if err := s.cache.Set(ctx, sessionKey, currentSession, 7*24*time.Hour); err != nil {
		s.logger.Error("failed to update session in cache", "user_id", userId, "error", err)
		return apperror.NewInternal("Terjadi Kesalahan saat memperbarui sesi.")
	}

	s.cache.Delete(ctx, otpKey)

	s.logger.Info("StepUp OTP verified successfully, session context updated", "user_id", userId, "new_ip", loginContext.IPAddress)

	return nil
}
