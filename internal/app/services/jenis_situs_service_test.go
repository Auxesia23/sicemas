package services_test

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"os"
	"testing"

	"sicemas/internal/app/services"
	"sicemas/internal/dto"
	"sicemas/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestJenisSitusService_CreateJenisSitus(t *testing.T) {
	defer func() {
		os.Unsetenv("AES_256_KEY")
		os.Unsetenv("PEPPER")
	}()
	os.Setenv("AES_256_KEY", "YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWE=")
	os.Setenv("PEPPER", "dummy-pepper")

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	ctx := context.Background()
	actorID := uuid.New()

	tests := []struct {
		name           string
		input         *dto.JenisSitusRequest
		mockSetup     func(repo *MockJenisSitusRepo, activityRepo *MockActivityRepo)
		expectedResp  *dto.JenisSitusResponse
		expectError    bool
	}{
		{
			name: "Success - CreateJenisSitus",
			input: &dto.JenisSitusRequest{
				Nama:      "Masjid",
				Deskripsi: "Tempat ibadah umat Islam",
			},
			mockSetup: func(repo *MockJenisSitusRepo, activityRepo *MockActivityRepo) {
				newJenisSitus := &entity.JenisSitus{
					ID:        uuid.New(),
					NamaJenis: "Masjid",
					Deskripsi: "Tempat ibadah umat Islam",
				}
				repo.On("Create", mock.Anything, mock.AnythingOfType("*dto.JenisSitusRequest")).
					Return(newJenisSitus, nil).Once()
				activityRepo.On("InsertActivity", mock.Anything, mock.AnythingOfType("*entity.Activity")).
					Return(nil).Once()
			},
			expectedResp: &dto.JenisSitusResponse{
				ID:          uuid.Nil,
				Nama:        "Masjid",
				Deskripsi:   "Tempat ibadah umat Islam",
				JumlahSitus: 0,
			},
			expectError: false,
		},
		{
			name: "Error - Repository Returns Error",
			input: &dto.JenisSitusRequest{
				Nama:      "Masjid",
				Deskripsi: "Tempat ibadah umat Islam",
			},
			mockSetup: func(repo *MockJenisSitusRepo, activityRepo *MockActivityRepo) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*dto.JenisSitusRequest")).
					Return(nil, errors.New("database error")).Once()
			},
			expectedResp: nil,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockJenisSitusRepo)
			activityRepo := new(MockActivityRepo)

			tt.mockSetup(repo, activityRepo)

			service := services.NewJenisSitusService(repo, activityRepo, logger)

			result, err := service.CreateJenisSitus(ctx, tt.input, actorID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp.Nama, result.Nama)
				assert.Equal(t, tt.expectedResp.Deskripsi, result.Deskripsi)
			}

			repo.AssertExpectations(t)
			activityRepo.AssertExpectations(t)
		})
	}
}

func TestJenisSitusService_GetAllJenisSitus(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	ctx := context.Background()

	tests := []struct {
		name           string
		mockSetup     func(repo *MockJenisSitusRepo)
		expectedResp []dto.JenisSitusResponse
		expectError  bool
	}{
		{
			name: "Success - Returns Data",
			mockSetup: func(repo *MockJenisSitusRepo) {
				data := []dto.JenisSitusResponse{
					{ID: uuid.New(), Nama: "Masjid", Deskripsi: "Tempat Ibadah", JumlahSitus: 10},
					{ID: uuid.New(), Nama: "Musholla", Deskripsi: "Tempat Ibadah Kecil", JumlahSitus: 5},
				}
				repo.On("ReadAll", mock.Anything).Return(data, nil).Once()
			},
			expectedResp: []dto.JenisSitusResponse{
				{ID: uuid.Nil, Nama: "Masjid", Deskripsi: "Tempat Ibadah", JumlahSitus: 10},
				{ID: uuid.Nil, Nama: "Musholla", Deskripsi: "Tempat Ibadah Kecil", JumlahSitus: 5},
			},
			expectError: false,
		},
		{
			name: "Success - Empty Data",
			mockSetup: func(repo *MockJenisSitusRepo) {
				repo.On("ReadAll", mock.Anything).Return([]dto.JenisSitusResponse{}, nil).Once()
			},
			expectedResp: []dto.JenisSitusResponse{},
			expectError:  false,
		},
		{
			name: "Error - Repository Returns Error",
			mockSetup: func(repo *MockJenisSitusRepo) {
				repo.On("ReadAll", mock.Anything).
					Return([]dto.JenisSitusResponse(nil), errors.New("database error")).Once()
			},
			expectedResp: nil,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockJenisSitusRepo)

			tt.mockSetup(repo)

			service := services.NewJenisSitusService(repo, nil, logger)

			result, err := service.GetAllJenisSitus(ctx)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedResp), len(result))
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestJenisSitusService_UpdateJenisSitus(t *testing.T) {
	defer func() {
		os.Unsetenv("AES_256_KEY")
		os.Unsetenv("PEPPER")
	}()
	os.Setenv("AES_256_KEY", "YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWE=")
	os.Setenv("PEPPER", "dummy-pepper")

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	ctx := context.Background()
	id := uuid.New()
	actorID := uuid.New()

	tests := []struct {
		name        string
		input       *dto.JenisSitusRequest
		mockSetup   func(repo *MockJenisSitusRepo, activityRepo *MockActivityRepo)
		expectError bool
	}{
		{
			name: "Success - UpdateJenisSitus",
			input: &dto.JenisSitusRequest{
				Nama:      "Masjid Updated",
				Deskripsi: "Tempat ibadah umat Islam yang sudah diperbarui",
			},
			mockSetup: func(repo *MockJenisSitusRepo, activityRepo *MockActivityRepo) {
				repo.On("Update", mock.Anything, id, mock.AnythingOfType("*dto.JenisSitusRequest")).
					Return(nil).Once()
				activityRepo.On("InsertActivity", mock.Anything, mock.AnythingOfType("*entity.Activity")).
					Return(nil).Once()
			},
			expectError: false,
		},
		{
			name: "Error - JenisSitus Not Found",
			input: &dto.JenisSitusRequest{
				Nama:      "Masjid Updated",
				Deskripsi: "Tempat ibadah umat Islam yang sudah diperbarui",
			},
			mockSetup: func(repo *MockJenisSitusRepo, activityRepo *MockActivityRepo) {
				repo.On("Update", mock.Anything, id, mock.AnythingOfType("*dto.JenisSitusRequest")).
					Return(sql.ErrNoRows).Once()
			},
			expectError: true,
		},
		{
			name: "Error - Repository Returns Error",
			input: &dto.JenisSitusRequest{
				Nama:      "Masjid Updated",
				Deskripsi: "Tempat ibadah umat Islam yang sudah diperbarui",
			},
			mockSetup: func(repo *MockJenisSitusRepo, activityRepo *MockActivityRepo) {
				repo.On("Update", mock.Anything, id, mock.AnythingOfType("*dto.JenisSitusRequest")).
					Return(errors.New("database error")).Once()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockJenisSitusRepo)
			activityRepo := new(MockActivityRepo)

			tt.mockSetup(repo, activityRepo)

			service := services.NewJenisSitusService(repo, activityRepo, logger)

			err := service.UpdateJenisSitus(ctx, id, tt.input, actorID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
			activityRepo.AssertExpectations(t)
		})
	}
}

func TestJenisSitusService_DeleteJenisSitus(t *testing.T) {
	defer func() {
		os.Unsetenv("AES_256_KEY")
		os.Unsetenv("PEPPER")
	}()
	os.Setenv("AES_256_KEY", "YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWE=")
	os.Setenv("PEPPER", "dummy-pepper")

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	ctx := context.Background()
	id := uuid.New()
	actorID := uuid.New()

	tests := []struct {
		name        string
		mockSetup   func(repo *MockJenisSitusRepo, activityRepo *MockActivityRepo)
		expectError bool
	}{
		{
			name: "Success - DeleteJenisSitus",
			mockSetup: func(repo *MockJenisSitusRepo, activityRepo *MockActivityRepo) {
				target := &dto.JenisSitusResponse{
					ID:        id,
					Nama:      "Masjid",
					Deskripsi: "Tempat Ibadah",
				}
				repo.On("ReadOne", mock.Anything, id).Return(target, nil).Once()
				repo.On("Delete", mock.Anything, id).Return(nil).Once()
				activityRepo.On("InsertActivity", mock.Anything, mock.AnythingOfType("*entity.Activity")).
					Return(nil).Once()
			},
			expectError: false,
		},
		{
			name: "Error - JenisSitus Not Found When Reading",
			mockSetup: func(repo *MockJenisSitusRepo, activityRepo *MockActivityRepo) {
				repo.On("ReadOne", mock.Anything, id).Return(nil, sql.ErrNoRows).Once()
			},
			expectError: true,
		},
		{
			name: "Error - Cannot Delete Due to Foreign Key Constraint",
			mockSetup: func(repo *MockJenisSitusRepo, activityRepo *MockActivityRepo) {
				target := &dto.JenisSitusResponse{
					ID:        id,
					Nama:      "Masjid",
					Deskripsi: "Tempat Ibadah",
				}
				repo.On("ReadOne", mock.Anything, id).Return(target, nil).Once()
				repo.On("Delete", mock.Anything, id).
					Return(errors.New("violates RESTRICT setting of foreign key constraint")).Once()
			},
			expectError: true,
		},
		{
			name: "Error - Not Found During Delete",
			mockSetup: func(repo *MockJenisSitusRepo, activityRepo *MockActivityRepo) {
				target := &dto.JenisSitusResponse{
					ID:        id,
					Nama:      "Masjid",
					Deskripsi: "Tempat Ibadah",
				}
				repo.On("ReadOne", mock.Anything, id).Return(target, nil).Once()
				repo.On("Delete", mock.Anything, id).Return(sql.ErrNoRows).Once()
			},
			expectError: true,
		},
		{
			name: "Error - Repository Delete Returns Generic Error",
			mockSetup: func(repo *MockJenisSitusRepo, activityRepo *MockActivityRepo) {
				target := &dto.JenisSitusResponse{
					ID:        id,
					Nama:      "Masjid",
					Deskripsi: "Tempat Ibadah",
				}
				repo.On("ReadOne", mock.Anything, id).Return(target, nil).Once()
				repo.On("Delete", mock.Anything, id).Return(errors.New("database error")).Once()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockJenisSitusRepo)
			activityRepo := new(MockActivityRepo)

			tt.mockSetup(repo, activityRepo)

			service := services.NewJenisSitusService(repo, activityRepo, logger)

			err := service.DeleteJenisSitus(ctx, id, actorID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
			activityRepo.AssertExpectations(t)
		})
	}
}