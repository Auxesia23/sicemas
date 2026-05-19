package services_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"sicemas/internal/app/services"
	"sicemas/internal/dto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupDashboardService() (*MockActivityRepo, *MockDashboardRepo, services.DashboardService) {
	mockActivityRepo := new(MockActivityRepo)
	mockDashboardRepo := new(MockDashboardRepo)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	svc := services.NewDashboardService(
		mockActivityRepo,
		mockDashboardRepo,
		logger,
	)

	return mockActivityRepo, mockDashboardRepo, svc
}

func TestDashboardServiceImpl_GetDashboardData(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		setupMock func(*MockDashboardRepo, *MockActivityRepo)
		expectErr bool
	}{
		{
			name: "Success - Get Dashboard Data",
			setupMock: func(md *MockDashboardRepo, ma *MockActivityRepo) {
				md.On("GetRecentActivities", mock.Anything).Return([]dto.ActivityResponse{
					{ID: uuid.New(), User: "Admin", Action: "Menambahkan petugas", Target: "Budi", Entity: "USER", CreatedAt: time.Now()},
				}, nil).Once()

				md.On("GetDashboardStats", mock.Anything).Return(&dto.Stats{
					TotalSitus:         100,
					TotalJenis:         10,
					PetugasAktif:       5,
					MenungguVerifikasi: 20,
				}, nil).Once()

				md.On("GetRecentSites", mock.Anything).Return([]dto.SitusKeagamaanResponse{
					{ID: uuid.New(), Nama: "Masjid Al-Falah", Jenis: "Masjid", Lokasi: "Desa Sukamaju", Status: "terverifikasi"},
				}, nil).Once()

				md.On("GetStatistikJenis", mock.Anything).Return([]dto.StatistikJenis{
					{Nama: "Masjid", JumlahSitus: 50},
					{Nama: "Musholla", JumlahSitus: 30},
				}, nil).Once()
			},
			expectErr: false,
		},
		{
			name: "Error - Get Recent Activities Failed",
			setupMock: func(md *MockDashboardRepo, ma *MockActivityRepo) {
				md.On("GetRecentActivities", mock.Anything).
					Return([]dto.ActivityResponse(nil), errors.New("database error")).Once()
			},
			expectErr: true,
		},
		{
			name: "Error - Get Dashboard Stats Failed",
			setupMock: func(md *MockDashboardRepo, ma *MockActivityRepo) {
				md.On("GetRecentActivities", mock.Anything).Return([]dto.ActivityResponse{
					{ID: uuid.New(), User: "Admin", Action: "Menambahkan petugas", Target: "Budi", Entity: "USER", CreatedAt: time.Now()},
				}, nil).Once()

				md.On("GetDashboardStats", mock.Anything).
					Return((*dto.Stats)(nil), errors.New("database error")).Once()
			},
			expectErr: true,
		},
		{
			name: "Error - Get Recent Sites Failed",
			setupMock: func(md *MockDashboardRepo, ma *MockActivityRepo) {
				md.On("GetRecentActivities", mock.Anything).Return([]dto.ActivityResponse{
					{ID: uuid.New(), User: "Admin", Action: "Menambahkan petugas", Target: "Budi", Entity: "USER", CreatedAt: time.Now()},
				}, nil).Once()

				md.On("GetDashboardStats", mock.Anything).Return(&dto.Stats{
					TotalSitus:         100,
					TotalJenis:         10,
					PetugasAktif:       5,
					MenungguVerifikasi: 20,
				}, nil).Once()

				md.On("GetRecentSites", mock.Anything).
					Return([]dto.SitusKeagamaanResponse(nil), errors.New("database error")).Once()
			},
			expectErr: true,
		},
		{
			name: "Error - Get Statistik Jenis Failed",
			setupMock: func(md *MockDashboardRepo, ma *MockActivityRepo) {
				md.On("GetRecentActivities", mock.Anything).Return([]dto.ActivityResponse{
					{ID: uuid.New(), User: "Admin", Action: "Menambahkan petugas", Target: "Budi", Entity: "USER", CreatedAt: time.Now()},
				}, nil).Once()

				md.On("GetDashboardStats", mock.Anything).Return(&dto.Stats{
					TotalSitus:         100,
					TotalJenis:         10,
					PetugasAktif:       5,
					MenungguVerifikasi: 20,
				}, nil).Once()

				md.On("GetRecentSites", mock.Anything).Return([]dto.SitusKeagamaanResponse{
					{ID: uuid.New(), Nama: "Masjid Al-Falah", Jenis: "Masjid", Lokasi: "Desa Sukamaju", Status: "terverifikasi"},
				}, nil).Once()

				md.On("GetStatistikJenis", mock.Anything).
					Return([]dto.StatistikJenis(nil), errors.New("database error")).Once()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockActivityRepo, mockDashboardRepo, svc := setupDashboardService()
			mockDashboardRepo.ExpectedCalls = nil

			tt.setupMock(mockDashboardRepo, mockActivityRepo)

			result, err := svc.GetDashboardData(ctx)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, result.Stats)
			}

			mockDashboardRepo.AssertExpectations(t)
		})
	}
}
