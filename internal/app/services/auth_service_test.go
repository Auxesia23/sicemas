package services_test

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"sicemas/internal/app/services"
	"sicemas/internal/dto"
	"sicemas/internal/entity"
	"sicemas/internal/utils"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAuthService() (*MockUserRepo, *MockCache, services.AuthService) {
	dummyKeyBytes := []byte("12345678901234567890123456789012")
	os.Setenv("AES_256_KEY", base64.StdEncoding.EncodeToString(dummyKeyBytes))
	os.Setenv("PEPPER", "dummy-pepper-rahasia")
	os.Setenv("JWT_SECRET", "test-jwt-secret-key-minimum-32-chars!")

	mockUserRepo := new(MockUserRepo)
	mockCache := new(MockCache)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	svc := services.NewAuthService(
		mockUserRepo,
		mockCache,
		logger,
	)

	return mockUserRepo, mockCache, svc
}

func TestAuthServiceImpl_Login(t *testing.T) {
	defer func() {
		os.Unsetenv("AES_256_KEY")
		os.Unsetenv("PEPPER")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("GOWA_URL")
	}()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()
	os.Setenv("GOWA_URL", mockServer.URL)
	os.Setenv("GOWA_DEVICE_ID", "test-device")

	ctx := context.Background()
	mockUserRepo, mockCache, svc := setupAuthService()

	tests := []struct {
		name      string
		input     *dto.UserLogin
		setupMock func(*MockUserRepo, *MockCache)
		expectErr bool
	}{
		{
			name: "Success - Login",
			input: &dto.UserLogin{
				NIP: "1234567890123456",
			},
			setupMock: func(mu *MockUserRepo, mc *MockCache) {
				encryptedNIP, _ := utils.Encrypt("1234567890123456")
				encryptedPhone, _ := utils.Encrypt("628123456789")
				user := &entity.User{
					ID:           uuid.New(),
					NIP:          encryptedNIP,
					NamaLengkap:  "Test User",
					NomorTelepon: encryptedPhone,
				}
				mu.On("ReadOne", mock.Anything, mock.Anything).Return(user, nil).Once()
				mc.On("Set", mock.Anything, mock.Anything, mock.Anything, 5*time.Minute).Return(nil).Once()
			},
			expectErr: false,
		},
		{
			name: "Error - User Not Found",
			input: &dto.UserLogin{
				NIP: "1234567890123456",
			},
			setupMock: func(mu *MockUserRepo, mc *MockCache) {
				mu.On("ReadOne", mock.Anything, mock.Anything).
					Return((*entity.User)(nil), sql.ErrNoRows).Once()
			},
			expectErr: true,
		},
		{
			name: "Error - Cache Set Failed",
			input: &dto.UserLogin{
				NIP: "1234567890123456",
			},
			setupMock: func(mu *MockUserRepo, mc *MockCache) {
				encryptedNIP, _ := utils.Encrypt("1234567890123456")
				encryptedPhone, _ := utils.Encrypt("628123456789")
				user := &entity.User{
					ID:           uuid.New(),
					NIP:          encryptedNIP,
					NamaLengkap:  "Test User",
					NomorTelepon: encryptedPhone,
				}
				mu.On("ReadOne", mock.Anything, mock.Anything).Return(user, nil).Once()
				mc.On("Set", mock.Anything, mock.Anything, mock.Anything, 5*time.Minute).
					Return(errors.New("redis error")).Once()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.ExpectedCalls = nil
			mockCache.ExpectedCalls = nil

			tt.setupMock(mockUserRepo, mockCache)

			err := svc.Login(ctx, tt.input)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockUserRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}

func TestAuthServiceImpl_VerifyOTP(t *testing.T) {
	defer func() {
		os.Unsetenv("AES_256_KEY")
		os.Unsetenv("PEPPER")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("GOWA_URL")
	}()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()
	os.Setenv("GOWA_URL", mockServer.URL)

	ctx := context.Background()
	mockUserRepo, mockCache, svc := setupAuthService()

	tests := []struct {
		name      string
		input     *dto.UserVerifyOTP
		loginCtx  *dto.SessionRequest
		setupMock func(*MockUserRepo, *MockCache)
		expectErr bool
	}{
		{
			name: "Success - Verify OTP",
			input: &dto.UserVerifyOTP{
				NIP: "1234567890123456",
				OTP: "123456",
			},
			loginCtx: &dto.SessionRequest{
				IPAddress: "127.0.0.1",
			},
			setupMock: func(mu *MockUserRepo, mc *MockCache) {
				encryptedNIP, _ := utils.Encrypt("1234567890123456")
				user := &entity.User{
					ID:          uuid.New(),
					NIP:         encryptedNIP,
					NamaLengkap: "Test User",
				}
				mu.On("ReadOne", mock.Anything, mock.Anything).Return(user, nil).Once()
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					ptr := args.Get(2).(*string)
					*ptr = "123456"
				}).Return(nil).Once()
				mc.On("Set", mock.Anything, mock.Anything, mock.Anything, 7*24*time.Hour).Return(nil).Once()
				mc.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectErr: false,
		},
		{
			name: "Error - User Not Found",
			input: &dto.UserVerifyOTP{
				NIP: "1234567890123456",
				OTP: "123456",
			},
			loginCtx: &dto.SessionRequest{
				IPAddress: "127.0.0.1",
			},
			setupMock: func(mu *MockUserRepo, mc *MockCache) {
				mu.On("ReadOne", mock.Anything, mock.Anything).
					Return((*entity.User)(nil), sql.ErrNoRows).Once()
			},
			expectErr: true,
		},
		{
			name: "Error - No OTP Found",
			input: &dto.UserVerifyOTP{
				NIP: "1234567890123456",
				OTP: "123456",
			},
			loginCtx: &dto.SessionRequest{
				IPAddress: "127.0.0.1",
			},
			setupMock: func(mu *MockUserRepo, mc *MockCache) {
				encryptedNIP, _ := utils.Encrypt("1234567890123456")
				user := &entity.User{
					ID:          uuid.New(),
					NIP:         encryptedNIP,
					NamaLengkap: "Test User",
				}
				mu.On("ReadOne", mock.Anything, mock.Anything).Return(user, nil).Once()
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).
					Return(redis.Nil).Once()
			},
			expectErr: true,
		},
		{
			name: "Error - Invalid OTP",
			input: &dto.UserVerifyOTP{
				NIP: "1234567890123456",
				OTP: "wrong",
			},
			loginCtx: &dto.SessionRequest{
				IPAddress: "127.0.0.1",
			},
			setupMock: func(mu *MockUserRepo, mc *MockCache) {
				encryptedNIP, _ := utils.Encrypt("1234567890123456")
				user := &entity.User{
					ID:          uuid.New(),
					NIP:         encryptedNIP,
					NamaLengkap: "Test User",
				}
				mu.On("ReadOne", mock.Anything, mock.Anything).Return(user, nil).Once()
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					ptr := args.Get(2).(*string)
					*ptr = "123456"
				}).Return(nil).Once()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.ExpectedCalls = nil
			mockCache.ExpectedCalls = nil

			tt.setupMock(mockUserRepo, mockCache)

			_, err := svc.VerifyOTP(ctx, tt.input, tt.loginCtx)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockUserRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}

func TestAuthServiceImpl_RefreshToken(t *testing.T) {
	defer func() {
		os.Unsetenv("AES_256_KEY")
		os.Unsetenv("PEPPER")
		os.Unsetenv("JWT_SECRET")
	}()

	ctx := context.Background()
	mockUserRepo, mockCache, svc := setupAuthService()

	validJTI := uuid.New()
	dummyUser := &entity.User{
		ID:          uuid.New(),
		NamaLengkap: "Test User",
	}
	validToken, _ := utils.GenerateRefreshToken(dummyUser, validJTI)

	tests := []struct {
		name         string
		refreshToken string
		requestCtx   *dto.SessionRequest
		setupMock    func(*MockUserRepo, *MockCache)
		expectErr    bool
	}{
		{
			name:         "Success - Refresh Token",
			refreshToken: validToken,
			requestCtx: &dto.SessionRequest{
				IPAddress: "127.0.0.1",
				UserAgent: "test-agent",
			},
			setupMock: func(mu *MockUserRepo, mc *MockCache) {
				mu.On("ReadById", mock.Anything, mock.Anything).Return(dummyUser, nil).Once()
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					ptr := args.Get(2).(*dto.SessionValue)
					*ptr = dto.SessionValue{
						UserID:    dummyUser.ID,
						IPAddress: "127.0.0.1",
						UserAgent: "test-agent",
					}
				}).Return(nil).Once()
				mc.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
				mc.On("Set", mock.Anything, mock.Anything, mock.Anything, 7*24*time.Hour).Return(nil).Once()
			},
			expectErr: false,
		},
		{
			name:         "Error - Invalid Token",
			refreshToken: "invalid.token.string",
			requestCtx: &dto.SessionRequest{
				IPAddress: "127.0.0.1",
			},
			setupMock: func(mu *MockUserRepo, mc *MockCache) {
			},
			expectErr: true,
		},
		{
			name:         "Error - User Not Found",
			refreshToken: validToken,
			requestCtx: &dto.SessionRequest{
				IPAddress: "127.0.0.1",
			},
			setupMock: func(mu *MockUserRepo, mc *MockCache) {
				mu.On("ReadById", mock.Anything, mock.Anything).
					Return((*entity.User)(nil), sql.ErrNoRows).Once()
			},
			expectErr: true,
		},
		{
			name:         "Error - Session Not Found",
			refreshToken: validToken,
			requestCtx: &dto.SessionRequest{
				IPAddress: "127.0.0.1",
			},
			setupMock: func(mu *MockUserRepo, mc *MockCache) {
				mu.On("ReadById", mock.Anything, mock.Anything).Return(dummyUser, nil).Once()
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).
					Return(redis.Nil).Once()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.ExpectedCalls = nil
			mockCache.ExpectedCalls = nil

			tt.setupMock(mockUserRepo, mockCache)

			_, err := svc.RefreshToken(ctx, tt.refreshToken, tt.requestCtx)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockUserRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}

func TestAuthServiceImpl_Logout(t *testing.T) {
	defer func() {
		os.Unsetenv("AES_256_KEY")
		os.Unsetenv("PEPPER")
		os.Unsetenv("JWT_SECRET")
	}()

	ctx := context.Background()
	mockUserRepo, mockCache, svc := setupAuthService()

	validJTI := uuid.New()
	dummyUser := &entity.User{
		ID:          uuid.New(),
		NamaLengkap: "Test User",
	}
	validToken, _ := utils.GenerateRefreshToken(dummyUser, validJTI)

	tests := []struct {
		name         string
		refreshToken string
		accessToken  *string
		setupMock    func(*MockUserRepo, *MockCache)
		expectErr    bool
	}{
		{
			name:         "Success - Logout",
			refreshToken: validToken,
			accessToken:  strPtr("valid.access.token"),
			setupMock: func(mu *MockUserRepo, mc *MockCache) {
				mc.On("Set", mock.Anything, mock.Anything, true, 15*time.Minute).Return(nil).Once()
				mc.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectErr: false,
		},
		{
			name:         "Success - Logout Without Access Token",
			refreshToken: validToken,
			accessToken:  nil,
			setupMock: func(mu *MockUserRepo, mc *MockCache) {
				mc.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectErr: false,
		},
		{
			name:         "Error - Invalid Refresh Token",
			refreshToken: "invalid.token",
			accessToken:  nil,
			setupMock: func(mu *MockUserRepo, mc *MockCache) {
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.ExpectedCalls = nil
			mockCache.ExpectedCalls = nil

			tt.setupMock(mockUserRepo, mockCache)

			err := svc.Logout(ctx, tt.refreshToken, tt.accessToken)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockUserRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}
func strPtr(s string) *string {
	return &s
}
