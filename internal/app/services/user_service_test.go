package services_test

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"sicemas/internal/app/services"
	"sicemas/internal/dto"
	"sicemas/internal/entity"
	"sicemas/internal/utils"

	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupUserService() (*MockUserRepo, *MockActivityRepo, *MockCache, services.UserService) {
	dummyKeyBytes := []byte("12345678901234567890123456789012")
	os.Setenv("AES_256_KEY", base64.StdEncoding.EncodeToString(dummyKeyBytes))
	os.Setenv("PEPPER", "dummy-pepper-rahasia")

	mockUserRepo := new(MockUserRepo)
	mockActivityRepo := new(MockActivityRepo)
	mockCache := new(MockCache)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	m, _ := casbinModel.NewModelFromString(`
		[request_definition]
		r = sub, obj, act
		[policy_definition]
		p = sub, obj, act
		[role_definition]
		g = _, _
		[policy_effect]
		e = some(where (p.eft == allow))
		[matchers]
		m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
	`)
	enforcer, _ := casbin.NewEnforcer(m, &DummyCasbinAdapter{})

	svc := services.NewUserService(
		mockUserRepo,
		mockActivityRepo,
		enforcer,
		mockCache,
		logger,
	)

	return mockUserRepo, mockActivityRepo, mockCache, svc
}

func TestUserServiceImpl_Register(t *testing.T) {
	defer os.Unsetenv("AES_256_KEY")
	defer os.Unsetenv("PEPPER")

	mockUser, mockAct, _, svc := setupUserService()

	actorId := uuid.New()
	unitKerja := "KUA Kecamatan Ciemas"
	ctx := context.Background()

	tests := []struct {
		name          string
		input        *dto.UserRequest
		mockSetup   func()
		expectedErr bool
	}{
		{
			name: "Success - User Registered",
			input: &dto.UserRequest{
				NamaLengkap:  "Budi Santoso",
				NIP:          "1234567890123456",
				Email:        "budi@test.com",
				NomorTelepon: "08123456789",
				UnitKerja:    &unitKerja,
				Peran:        "ADMIN",
				Jabatan:      "Staff",
			},
			mockSetup: func() {
				mockUser.On("Create", mock.Anything, mock.AnythingOfType("*dto.UserRequest"), mock.Anything).
					Return(uuid.New().String(), nil).Once()
				mockAct.On("InsertActivity", mock.Anything, mock.AnythingOfType("*entity.Activity")).
					Return(nil).Once()
			},
			expectedErr: false,
		},
		{
			name: "Error - Duplicate NIP",
			input: &dto.UserRequest{
				NamaLengkap:  "Budi Santoso",
				NIP:          "1234567890123456",
				Email:        "budi@test.com",
				NomorTelepon: "08123456789",
				UnitKerja:    &unitKerja,
				Peran:        "ADMIN",
				Jabatan:      "Staff",
			},
			mockSetup: func() {
				mockUser.On("Create", mock.Anything, mock.Anything, mock.Anything).
					Return("", errors.New("duplicate key value violates unique constraint")).Once()
			},
			expectedErr: true,
		},
		{
			name: "Error - Database Create Failed",
			input: &dto.UserRequest{
				NamaLengkap:  "Budi Santoso",
				NIP:          "1234567890123456",
				Email:        "budi@test.com",
				NomorTelepon: "08123456789",
				UnitKerja:    &unitKerja,
				Peran:        "ADMIN",
				Jabatan:      "Staff",
			},
			mockSetup: func() {
				mockUser.On("Create", mock.Anything, mock.Anything, mock.Anything).
					Return("", errors.New("database connection failed")).Once()
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser.ExpectedCalls = nil
			mockAct.ExpectedCalls = nil

			tt.mockSetup()

			reqCopy := *tt.input

			err := svc.Register(ctx, &reqCopy, actorId)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockUser.AssertExpectations(t)
			mockAct.AssertExpectations(t)
		})
	}
}

func TestUserServiceImpl_GetAllUser(t *testing.T) {
	defer os.Unsetenv("AES_256_KEY")
	defer os.Unsetenv("PEPPER")

	mockUser, _, mockCache, svc := setupUserService()

	ctx := context.Background()

	tests := []struct {
		name          string
		mockSetup    func()
		expectedErr  bool
		expectedCnt int
	}{
		{
			name: "Success - Returns Data",
			mockSetup: func() {
				encNIP1, _ := utils.Encrypt("198001012010011001")
				encEmail1, _ := utils.Encrypt("user1@test.com")
				encPhone1, _ := utils.Encrypt("08123456789")

				encNIP2, _ := utils.Encrypt("198001012010011002")
				encEmail2, _ := utils.Encrypt("user2@test.com")
				encPhone2, _ := utils.Encrypt("08123456780")

				users := []entity.User{
					{
						ID:           uuid.New(),
						NamaLengkap:  "User One",
						Jabatan:      "Staff",
						UnitKerja:    "KUA Ciemas",
						NIP:          encNIP1,
						Email:        encEmail1,
						NomorTelepon: encPhone1,
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
					},
					{
						ID:           uuid.New(),
						NamaLengkap:  "User Two",
						Jabatan:      "Staff",
						UnitKerja:    "KUA Ciemas",
						NIP:          encNIP2,
						Email:        encEmail2,
						NomorTelepon: encPhone2,
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
					},
				}
				mockUser.On("ReadAll", mock.Anything).Return(users, nil).Once()
			},
			expectedErr:  false,
			expectedCnt: 2,
		},
		{
			name: "Success - Empty Data",
			mockSetup: func() {
				mockUser.On("ReadAll", mock.Anything).Return([]entity.User{}, nil).Once()
			},
			expectedErr:  false,
			expectedCnt: 0,
		},
		{
			name: "Error - Repository Returns Error",
			mockSetup: func() {
				mockUser.On("ReadAll", mock.Anything).
					Return([]entity.User(nil), errors.New("database error")).Once()
			},
			expectedErr:  true,
			expectedCnt: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser.ExpectedCalls = nil

			tt.mockSetup()

			res, err := svc.GetAllUser(ctx)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCnt, len(res))
			}

			mockUser.AssertExpectations(t)
		})
	}

	_ = mockCache
}

func TestUserServiceImpl_UpdateUser(t *testing.T) {
	defer os.Unsetenv("AES_256_KEY")
	defer os.Unsetenv("PEPPER")

	ctx := context.Background()
	userId := uuid.New()
	actorId := uuid.New()

	tests := []struct {
		name          string
		input        *dto.UserRequest
		setupMock    func(*MockUserRepo, *MockActivityRepo, services.UserService)
		expectedErr bool
	}{
		{
			name: "Success - Update User",
			input: &dto.UserRequest{
				NamaLengkap:  "Budi Updated",
				NIP:          "1234567890123456",
				Email:        "budi@test.com",
				NomorTelepon: "08123456789",
				Peran:        "ADMIN",
				Jabatan:      "Staff",
			},
			setupMock: func(mu *MockUserRepo, ma *MockActivityRepo, svc services.UserService) {
				mu.On("Update", mock.Anything, userId, mock.Anything, mock.Anything).
					Return(nil).Once()
				ma.On("InsertActivity", mock.Anything, mock.AnythingOfType("*entity.Activity")).
					Return(nil).Once()
			},
			expectedErr: false,
		},
		{
			name: "Error - User Not Found",
			input: &dto.UserRequest{
				NamaLengkap:  "Budi Updated",
				NIP:          "1234567890123456",
				Email:        "budi@test.com",
				NomorTelepon: "08123456789",
				Peran:        "ADMIN",
				Jabatan:      "Staff",
			},
			setupMock: func(mu *MockUserRepo, ma *MockActivityRepo, svc services.UserService) {
				mu.On("Update", mock.Anything, userId, mock.Anything, mock.Anything).
					Return(sql.ErrNoRows).Once()
			},
			expectedErr: true,
		},
		{
			name: "Error - Duplicate NIP",
			input: &dto.UserRequest{
				NamaLengkap:  "Budi Updated",
				NIP:          "1234567890123456",
				Email:        "budi@test.com",
				NomorTelepon: "08123456789",
				Peran:        "ADMIN",
				Jabatan:      "Staff",
			},
			setupMock: func(mu *MockUserRepo, ma *MockActivityRepo, svc services.UserService) {
				mu.On("Update", mock.Anything, userId, mock.Anything, mock.Anything).
					Return(errors.New("duplicate key value violates unique constraint")).Once()
			},
			expectedErr: true,
		},
		{
			name: "Error - Database Error",
			input: &dto.UserRequest{
				NamaLengkap:  "Budi Updated",
				NIP:          "1234567890123456",
				Email:        "budi@test.com",
				NomorTelepon: "08123456789",
				Peran:        "ADMIN",
				Jabatan:      "Staff",
			},
			setupMock: func(mu *MockUserRepo, ma *MockActivityRepo, svc services.UserService) {
				mu.On("Update", mock.Anything, userId, mock.Anything, mock.Anything).
					Return(errors.New("database error")).Once()
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser, mockAct, mockCache, svc := setupUserService()
			_ = mockCache
			
			tt.setupMock(mockUser, mockAct, svc)

			inputCopy := *tt.input

			err := svc.UpdateUser(ctx, userId, actorId, &inputCopy)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockUser.AssertExpectations(t)
			mockAct.AssertExpectations(t)
})
	}
}

func TestUserServiceImpl_DeleteUser(t *testing.T) {
	defer os.Unsetenv("AES_256_KEY")
	defer os.Unsetenv("PEPPER")

	mockUser, mockAct, mockCache, svc := setupUserService()

	ctx := context.Background()
	userId := uuid.New()
	actorId := uuid.New()

	tests := []struct {
		name          string
		mockSetup    func()
		expectedErr bool
	}{
		{
			name: "Success - Delete User",
			mockSetup: func() {
				encNIP, _ := utils.Encrypt("198001012010011001")
				encEmail, _ := utils.Encrypt("user@test.com")
				encPhone, _ := utils.Encrypt("08123456789")

				user := &entity.User{
					ID:           userId,
					NamaLengkap:  "User To Delete",
					Jabatan:      "Staff",
					UnitKerja:    "KUA Ciemas",
					NIP:          encNIP,
					Email:        encEmail,
					NomorTelepon: encPhone,
					CreatedAt:    time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockUser.On("ReadById", mock.Anything, userId.String()).Return(user, nil).Once()
				mockUser.On("Delete", mock.Anything, userId).Return(nil).Once()
				mockAct.On("InsertActivity", mock.Anything, mock.AnythingOfType("*entity.Activity")).
					Return(nil).Once()
			},
			expectedErr: false,
		},
		{
			name: "Error - User Not Found",
			mockSetup: func() {
				mockUser.On("ReadById", mock.Anything, userId.String()).
					Return((*entity.User)(nil), sql.ErrNoRows).Once()
			},
			expectedErr: true,
		},
		{
			name: "Error - Delete Failed",
			mockSetup: func() {
				encNIP, _ := utils.Encrypt("198001012010011001")
				encEmail, _ := utils.Encrypt("user@test.com")
				encPhone, _ := utils.Encrypt("08123456789")

				user := &entity.User{
					ID:           userId,
					NamaLengkap:  "User To Delete",
					Jabatan:      "Staff",
					UnitKerja:    "KUA Ciemas",
					NIP:          encNIP,
					Email:        encEmail,
					NomorTelepon: encPhone,
					CreatedAt:    time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockUser.On("ReadById", mock.Anything, userId.String()).Return(user, nil).Once()
				mockUser.On("Delete", mock.Anything, userId).
					Return(errors.New("database error")).Once()
			},
			expectedErr: true,
		},
		{
			name: "Error - User Not Found During Delete",
			mockSetup: func() {
				encNIP, _ := utils.Encrypt("198001012010011001")
				encEmail, _ := utils.Encrypt("user@test.com")
				encPhone, _ := utils.Encrypt("08123456789")

				user := &entity.User{
					ID:           userId,
					NamaLengkap:  "User To Delete",
					Jabatan:      "Staff",
					UnitKerja:    "KUA Ciemas",
					NIP:          encNIP,
					Email:        encEmail,
					NomorTelepon: encPhone,
					CreatedAt:    time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockUser.On("ReadById", mock.Anything, userId.String()).Return(user, nil).Once()
				mockUser.On("Delete", mock.Anything, userId).Return(sql.ErrNoRows).Once()
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser.ExpectedCalls = nil
			mockAct.ExpectedCalls = nil

			tt.mockSetup()

			err := svc.DeleteUser(ctx, userId, actorId)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockUser.AssertExpectations(t)
			mockAct.AssertExpectations(t)
		})
	}

	_ = mockCache
}

func TestUserServiceImpl_GetProfile(t *testing.T) {
	defer os.Unsetenv("AES_256_KEY")
	defer os.Unsetenv("PEPPER")

	mockUser, _, mockCache, svc := setupUserService()

	ctx := context.Background()
	userId := uuid.New()

	tests := []struct {
		name          string
		userID       string
		mockSetup   func()
		expectedErr bool
		expectedNm  string
		expectedNIP string
	}{
		{
			name:   "Success - Get User Profile",
			userID: userId.String(),
			mockSetup: func() {
				encNIP, _ := utils.Encrypt("198001012010011001")
				encEmail, _ := utils.Encrypt("siti@test.com")
				encPhone, _ := utils.Encrypt("08123456789")

				dummyUser := &entity.User{
					ID:           userId,
					NamaLengkap:  "Siti Aminah",
					Jabatan:      "Staff",
					UnitKerja:    "KUA Ciemas",
					NIP:          encNIP,
					Email:        encEmail,
					NomorTelepon: encPhone,
					CreatedAt:    time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockUser.On("ReadById", mock.Anything, userId.String()).
					Return(dummyUser, nil).Once()
			},
			expectedErr: false,
			expectedNm:  "Siti Aminah",
			expectedNIP: "198001012010011001",
		},
		{
			name:   "Error - User Not Found",
			userID: userId.String(),
			mockSetup: func() {
				mockUser.On("ReadById", mock.Anything, userId.String()).
					Return((*entity.User)(nil), sql.ErrNoRows).Once()
			},
			expectedErr: true,
			expectedNm:  "",
			expectedNIP: "",
		},
		{
			name:   "Error - Repository Error",
			userID: userId.String(),
			mockSetup: func() {
				mockUser.On("ReadById", mock.Anything, userId.String()).
					Return((*entity.User)(nil), errors.New("database error")).Once()
			},
			expectedErr: true,
			expectedNm:  "",
			expectedNIP: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser.ExpectedCalls = nil

			tt.mockSetup()

			res, err := svc.GetProfile(ctx, tt.userID)

			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tt.expectedNm, res.NamaLengkap)
				assert.Equal(t, tt.expectedNIP, res.NIP)
			}

			mockUser.AssertExpectations(t)
		})
	}

	_ = mockCache
}

func TestUserServiceImpl_GetPermissions(t *testing.T) {
	defer os.Unsetenv("AES_256_KEY")
	defer os.Unsetenv("PEPPER")

	mockUser, _, mockCache, svc := setupUserService()

	ctx := context.Background()
	userId := uuid.New()

	tests := []struct {
		name          string
		userID       string
		mockSetup   func()
		expectedErr  bool
		expectedPrs int
	}{
		{
			name:   "Success - Get User Permissions",
			userID: userId.String(),
			mockSetup: func() {
				encNIP, _ := utils.Encrypt("198001012010011001")
				encEmail, _ := utils.Encrypt("user@test.com")
				encPhone, _ := utils.Encrypt("08123456789")

				user := &entity.User{
					ID:           userId,
					NamaLengkap:  "User with Permissions",
					Jabatan:      "Staff",
					UnitKerja:    "KUA Ciemas",
					NIP:          encNIP,
					Email:        encEmail,
					NomorTelepon: encPhone,
					CreatedAt:    time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockUser.On("ReadById", mock.Anything, userId.String()).Return(user, nil).Once()
			},
			expectedErr:  false,
			expectedPrs: 0,
		},
		{
			name:   "Error - User Not Found",
			userID: userId.String(),
			mockSetup: func() {
				mockUser.On("ReadById", mock.Anything, userId.String()).
					Return((*entity.User)(nil), sql.ErrNoRows).Once()
			},
			expectedErr:  true,
			expectedPrs: 0,
		},
		{
			name:   "Error - Repository Error",
			userID: userId.String(),
			mockSetup: func() {
				mockUser.On("ReadById", mock.Anything, userId.String()).
					Return((*entity.User)(nil), errors.New("database error")).Once()
			},
			expectedErr:  true,
			expectedPrs: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser.ExpectedCalls = nil

			tt.mockSetup()

			res, err := svc.GetPermissions(ctx, tt.userID)

			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tt.expectedPrs, len(res.Permissions))
			}

			mockUser.AssertExpectations(t)
		})
	}

	_ = mockCache
}