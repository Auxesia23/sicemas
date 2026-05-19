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

	"sicemas/internal/app/services"
	"sicemas/internal/dto"
	"sicemas/internal/utils"

	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupSitusService() (*MockSitusKeagamaanRepo, *MockFotoSitusRepo, *MockActivityRepo, *casbin.Enforcer, services.SitusKeagamaanService) {
	dummyKeyBytes := []byte("12345678901234567890123456789012")
	os.Setenv("AES_256_KEY", base64.StdEncoding.EncodeToString(dummyKeyBytes))
	os.Setenv("PEPPER", "dummy-pepper-rahasia")

	mockSitusRepo := new(MockSitusKeagamaanRepo)
	mockFotoRepo := new(MockFotoSitusRepo)
	mockActRepo := new(MockActivityRepo)

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

	cld, _ := cloudinary.NewFromURL("cloudinary://123456789012345:dummy_secret@dummy_cloud")

	svc := services.NewSitusKeagamaanService(
		mockSitusRepo,
		mockFotoRepo,
		mockActRepo,
		cld,
		enforcer,
		logger,
	)

	return mockSitusRepo, mockFotoRepo, mockActRepo, enforcer, svc
}

func TestSitusKeagamaanServiceImpl_CreateSitusKeagamaan(t *testing.T) {
	defer func() {
		os.Unsetenv("AES_256_KEY")
		os.Unsetenv("PEPPER")
	}()

	ctx := context.Background()
	authorId := uuid.New()
	actorId := uuid.New()
	expectedSitusId := uuid.New()

	tests := []struct {
		name      string
		input     *dto.SitusKeagamaanRequest
		setupMock func(*MockSitusKeagamaanRepo, *MockActivityRepo)
		expectErr bool
	}{
		{
			name: "Success - Create Situs",
			input: &dto.SitusKeagamaanRequest{
				Nama:                 "Masjid Al-Ikhlas",
				NomorBadanHukum:      "12345",
				NomorAIW:             "AIW-123",
				NomorSertifikatWakaf: "WKF-123",
				NomorTelepon:         "08123456789",
				Email:                "masjid@test.com",
				NomorTelponPengurus:  []string{"08987654321"},
			},
			setupMock: func(ms *MockSitusKeagamaanRepo, ma *MockActivityRepo) {
				ms.On("Create", ctx, mock.AnythingOfType("*dto.SitusKeagamaanRequest"), authorId).
					Return(expectedSitusId, nil).Once()
				ma.On("InsertActivity", ctx, mock.AnythingOfType("*entity.Activity")).
					Return(nil).Once()
			},
			expectErr: false,
		},
		{
			name: "Error - Repository Failed",
			input: &dto.SitusKeagamaanRequest{
				Nama: "Gereja Katedral",
			},
			setupMock: func(ms *MockSitusKeagamaanRepo, ma *MockActivityRepo) {
				ms.On("Create", ctx, mock.AnythingOfType("*dto.SitusKeagamaanRequest"), authorId).
					Return(uuid.Nil, errors.New("database error")).Once()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSitus, _, mockAct, _, svc := setupSitusService()
			tt.setupMock(mockSitus, mockAct)

			reqCopy := *tt.input
			reqCopy.NomorTelponPengurus = append([]string(nil), tt.input.NomorTelponPengurus...)

			id, err := svc.CreateSitusKeagamaan(ctx, &reqCopy, authorId, actorId)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, uuid.Nil, id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedSitusId, id)
			}

			mockSitus.AssertExpectations(t)
			mockAct.AssertExpectations(t)
		})
	}
}

func TestSitusKeagamaanServiceImpl_GetAllSitusKeagamaan(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()

	tests := []struct {
		name      string
		setupMock func(*MockSitusKeagamaanRepo, *casbin.Enforcer)
		expectErr bool
		expectCnt int
	}{
		{
			name: "Success - Read All (Admin Role)",
			setupMock: func(ms *MockSitusKeagamaanRepo, enf *casbin.Enforcer) {
				enf.AddPolicy(userId.String(), "situs", "read_all")
				ms.On("ReadAll", ctx).Return([]dto.SitusKeagamaanResponse{{}, {}}, nil).Once()
			},
			expectErr: false,
			expectCnt: 2,
		},
		{
			name: "Success - Read Own (User Role)",
			setupMock: func(ms *MockSitusKeagamaanRepo, enf *casbin.Enforcer) {
				ms.On("ReadOwn", ctx, userId).Return([]dto.SitusKeagamaanResponse{{}}, nil).Once()
			},
			expectErr: false,
			expectCnt: 1,
		},
		{
			name: "Error - Repository Failed",
			setupMock: func(ms *MockSitusKeagamaanRepo, enf *casbin.Enforcer) {
				ms.On("ReadOwn", ctx, userId).Return([]dto.SitusKeagamaanResponse(nil), errors.New("db error")).Once()
			},
			expectErr: true,
			expectCnt: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSitus, _, _, enf, svc := setupSitusService()
			tt.setupMock(mockSitus, enf)

			res, err := svc.GetAllSitusKeagamaan(ctx, userId)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Len(t, res, tt.expectCnt)
			}

			mockSitus.AssertExpectations(t)
		})
	}
}

func TestSitusKeagamaanServiceImpl_GetDetailSitusKeagamaan(t *testing.T) {
	defer func() {
		os.Unsetenv("AES_256_KEY")
		os.Unsetenv("PEPPER")
	}()

	ctx := context.Background()
	situsId := uuid.New()
	userId := uuid.New()

	encPhone, _ := utils.Encrypt("08123456789")
	encEmail, _ := utils.Encrypt("test@email.com")

	tests := []struct {
		name      string
		setupMock func(*MockSitusKeagamaanRepo, *MockFotoSitusRepo, *casbin.Enforcer)
		expectErr bool
	}{
		{
			name: "Success - Get Detail (Admin)",
			setupMock: func(ms *MockSitusKeagamaanRepo, mf *MockFotoSitusRepo, enf *casbin.Enforcer) {
				enf.AddPolicy(userId.String(), "situs", "read_all")

				dummyDetail := &dto.SitusKeagamaanDetailResponse{
					ID:                   situsId,
					NomorTelepon:         encPhone,
					Email:                encEmail,
					NomorBadanHukum:      encPhone,
					NomorAIW:             encPhone,
					NomorSertifikatWakaf: encPhone,
					NomorTelponPengurus:  []string{encPhone},
				}

				ms.On("ReadDetail", ctx, situsId).Return(dummyDetail, nil).Once()
				mf.On("GetBySitusID", ctx, situsId).Return(&[]dto.FotoResponse{{}}, nil).Once()
			},
			expectErr: false,
		},
		{
			name: "Error - Forbidden Ownership (User)",
			setupMock: func(ms *MockSitusKeagamaanRepo, mf *MockFotoSitusRepo, enf *casbin.Enforcer) {
				ms.On("CheckOwnership", ctx, situsId, userId).Return(sql.ErrNoRows).Once()
			},
			expectErr: true,
		},
		{
			name: "Error - Situs Not Found",
			setupMock: func(ms *MockSitusKeagamaanRepo, mf *MockFotoSitusRepo, enf *casbin.Enforcer) {
				enf.AddPolicy(userId.String(), "situs", "read_all")
				ms.On("ReadDetail", ctx, situsId).Return((*dto.SitusKeagamaanDetailResponse)(nil), sql.ErrNoRows).Once()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSitus, mockFoto, _, enf, svc := setupSitusService()
			tt.setupMock(mockSitus, mockFoto, enf)

			res, err := svc.GetDetailSitusKeagamaan(ctx, situsId, userId)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, "08123456789", res.NomorTelepon)
			}

			mockSitus.AssertExpectations(t)
			mockFoto.AssertExpectations(t)
		})
	}
}

func TestSitusKeagamaanServiceImpl_UpdateSitus(t *testing.T) {
	defer func() {
		os.Unsetenv("AES_256_KEY")
		os.Unsetenv("PEPPER")
	}()

	ctx := context.Background()
	situsId := uuid.New()
	userId := uuid.New()
	actorId := uuid.New()

	tests := []struct {
		name      string
		input     *dto.SitusKeagamaanUpdate
		setupMock func(*MockSitusKeagamaanRepo, *MockActivityRepo)
		expectErr bool
	}{
		{
			name: "Success - Update Situs",
			input: &dto.SitusKeagamaanUpdate{
				Nama:                 "Masjid Baru",
				NomorBadanHukum:      "123",
				NomorAIW:             "AIW-1",
				NomorSertifikatWakaf: "WKF-1",
				NomorTelepon:         "08111",
				Email:                "update@test.com",
				NomorTelponPengurus:  []string{"08222"},
			},
			setupMock: func(ms *MockSitusKeagamaanRepo, ma *MockActivityRepo) {
				ms.On("CheckOwnership", ctx, situsId, userId).Return(nil).Once()
				ms.On("Update", ctx, situsId, mock.AnythingOfType("*dto.SitusKeagamaanUpdate")).Return(nil).Once()
				ma.On("InsertActivity", ctx, mock.AnythingOfType("*entity.Activity")).Return(nil).Once()
			},
			expectErr: false,
		},
		{
			name: "Error - Forbidden Ownership",
			input: &dto.SitusKeagamaanUpdate{
				Nama: "Masjid Baru",
			},
			setupMock: func(ms *MockSitusKeagamaanRepo, ma *MockActivityRepo) {
				ms.On("CheckOwnership", ctx, situsId, userId).Return(sql.ErrNoRows).Once()
			},
			expectErr: true,
		},
		{
			name: "Error - Situs Not Found During Update",
			input: &dto.SitusKeagamaanUpdate{
				Nama:                 "Masjid Baru",
				NomorBadanHukum:      "123",
				NomorAIW:             "AIW-1",
				NomorSertifikatWakaf: "WKF-1",
				NomorTelepon:         "08111",
				Email:                "update@test.com",
				NomorTelponPengurus:  []string{"08222"},
			},
			setupMock: func(ms *MockSitusKeagamaanRepo, ma *MockActivityRepo) {
				ms.On("CheckOwnership", ctx, situsId, userId).Return(nil).Once()
				ms.On("Update", ctx, situsId, mock.AnythingOfType("*dto.SitusKeagamaanUpdate")).Return(sql.ErrNoRows).Once()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSitus, _, mockAct, _, svc := setupSitusService()
			tt.setupMock(mockSitus, mockAct)

			reqCopy := *tt.input
			reqCopy.NomorTelponPengurus = append([]string(nil), tt.input.NomorTelponPengurus...)

			err := svc.UpdateSitus(ctx, situsId, userId, &reqCopy, actorId)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockSitus.AssertExpectations(t)
			mockAct.AssertExpectations(t)
		})
	}
}

func TestSitusKeagamaanServiceImpl_DeleteSitus(t *testing.T) {
	ctx := context.Background()
	situsId := uuid.New()
	actorId := uuid.New()

	tests := []struct {
		name      string
		setupMock func(*MockSitusKeagamaanRepo, *MockActivityRepo)
		expectErr bool
	}{
		{
			name: "Success - Delete Situs",
			setupMock: func(ms *MockSitusKeagamaanRepo, ma *MockActivityRepo) {
				dummyDetail := &dto.SitusKeagamaanDetailResponse{Nama: "Masjid Test"}
				ms.On("ReadDetail", ctx, situsId).Return(dummyDetail, nil).Once()
				ms.On("Delete", ctx, situsId).Return(nil).Once()
				ma.On("InsertActivity", ctx, mock.AnythingOfType("*entity.Activity")).Return(nil).Once()
			},
			expectErr: false,
		},
		{
			name: "Error - Situs Not Found Before Delete",
			setupMock: func(ms *MockSitusKeagamaanRepo, ma *MockActivityRepo) {
				ms.On("ReadDetail", ctx, situsId).Return((*dto.SitusKeagamaanDetailResponse)(nil), sql.ErrNoRows).Once()
			},
			expectErr: true,
		},
		{
			name: "Error - Delete Failed",
			setupMock: func(ms *MockSitusKeagamaanRepo, ma *MockActivityRepo) {
				dummyDetail := &dto.SitusKeagamaanDetailResponse{Nama: "Masjid Test"}
				ms.On("ReadDetail", ctx, situsId).Return(dummyDetail, nil).Once()
				ms.On("Delete", ctx, situsId).Return(errors.New("db error")).Once()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSitus, _, mockAct, _, svc := setupSitusService()
			tt.setupMock(mockSitus, mockAct)

			err := svc.DeleteSitus(ctx, situsId, actorId)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockSitus.AssertExpectations(t)
			mockAct.AssertExpectations(t)
		})
	}
}

func TestSitusKeagamaanServiceImpl_VerifySitus(t *testing.T) {
	ctx := context.Background()
	situsId := uuid.New()
	actorId := uuid.New()
	req := &dto.VerifikasiSitusRequest{
		StatusVerifikasi: "terverifikasi",
		SitusID:          "SITUS-001",
	}

	tests := []struct {
		name      string
		setupMock func(*MockSitusKeagamaanRepo, *MockActivityRepo)
		expectErr bool
	}{
		{
			name: "Success - Verify Situs",
			setupMock: func(ms *MockSitusKeagamaanRepo, ma *MockActivityRepo) {
				dummyDetail := &dto.SitusKeagamaanDetailResponse{Nama: "Masjid Test"}
				ms.On("ReadDetail", ctx, situsId).Return(dummyDetail, nil).Once()
				ms.On("Verify", ctx, situsId, req).Return(nil).Once()
				ma.On("InsertActivity", ctx, mock.AnythingOfType("*entity.Activity")).Return(nil).Once()
			},
			expectErr: false,
		},
		{
			name: "Error - Situs Not Found",
			setupMock: func(ms *MockSitusKeagamaanRepo, ma *MockActivityRepo) {
				ms.On("ReadDetail", ctx, situsId).Return((*dto.SitusKeagamaanDetailResponse)(nil), sql.ErrNoRows).Once()
			},
			expectErr: true,
		},
		{
			name: "Error - Verify Failed",
			setupMock: func(ms *MockSitusKeagamaanRepo, ma *MockActivityRepo) {
				dummyDetail := &dto.SitusKeagamaanDetailResponse{Nama: "Masjid Test"}
				ms.On("ReadDetail", ctx, situsId).Return(dummyDetail, nil).Once()
				ms.On("Verify", ctx, situsId, req).Return(errors.New("db error")).Once()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSitus, _, mockAct, _, svc := setupSitusService()
			tt.setupMock(mockSitus, mockAct)

			err := svc.VerifySitus(ctx, situsId, req, actorId)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockSitus.AssertExpectations(t)
			mockAct.AssertExpectations(t)
		})
	}
}

func TestSitusKeagamaanServiceImpl_ExportSitusToExcel(t *testing.T) {
	defer func() {
		os.Unsetenv("AES_256_KEY")
		os.Unsetenv("PEPPER")
	}()

	ctx := context.Background()
	jenisSitus := "Masjid"
	actorId := uuid.New()

	// Setup dummy encrypted string biar utils.Decrypt nggak panic
	encDummy, _ := utils.Encrypt("dummy_data")

	tests := []struct {
		name      string
		setupMock func(*MockSitusKeagamaanRepo, *MockActivityRepo)
		expectErr bool
	}{
		{
			name: "Success - Export Excel",
			setupMock: func(ms *MockSitusKeagamaanRepo, ma *MockActivityRepo) {
				dummyList := []dto.SitusKeagamaanDetailResponse{
					{
						ID:                   uuid.New(),
						NomorTelepon:         encDummy,
						Email:                encDummy,
						NomorBadanHukum:      encDummy,
						NomorAIW:             encDummy,
						NomorSertifikatWakaf: encDummy,
						NomorTelponPengurus:  []string{encDummy},
					},
				}
				ms.On("ReadAllDetail", ctx, jenisSitus).Return(dummyList, nil).Once()
				ma.On("InsertActivity", ctx, mock.AnythingOfType("*entity.Activity")).Return(nil).Once()
			},
			expectErr: false,
		},
		{
			name: "Error - Repo ReadAllDetail Failed",
			setupMock: func(ms *MockSitusKeagamaanRepo, ma *MockActivityRepo) {
				ms.On("ReadAllDetail", ctx, jenisSitus).Return([]dto.SitusKeagamaanDetailResponse(nil), errors.New("db err")).Once()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSitus, _, mockAct, _, svc := setupSitusService()
			tt.setupMock(mockSitus, mockAct)

			file, err := svc.ExportSitusToExcel(ctx, jenisSitus, actorId)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, file)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, file)
			}

			mockSitus.AssertExpectations(t)
			mockAct.AssertExpectations(t)
		})
	}
}

func TestSitusKeagamaanServiceImpl_GetAllSitusForPublic(t *testing.T) {
	ctx := context.Background()
	filter := dto.PublicListFilter{}

	tests := []struct {
		name      string
		setupMock func(*MockSitusKeagamaanRepo)
		expectErr bool
		expectCnt int
	}{
		{
			name: "Success - Get Public List",
			setupMock: func(ms *MockSitusKeagamaanRepo) {
				ms.On("ReadForPublic", ctx, filter).Return([]dto.SitusPublicListResponse{{}}, nil).Once()
			},
			expectErr: false,
			expectCnt: 1,
		},
		{
			name: "Success - Empty Public List",
			setupMock: func(ms *MockSitusKeagamaanRepo) {
				ms.On("ReadForPublic", ctx, filter).Return([]dto.SitusPublicListResponse{}, nil).Once()
			},
			expectErr: false,
			expectCnt: 0,
		},
		{
			name: "Error - Repo Failed",
			setupMock: func(ms *MockSitusKeagamaanRepo) {
				ms.On("ReadForPublic", ctx, filter).Return([]dto.SitusPublicListResponse(nil), errors.New("db error")).Once()
			},
			expectErr: true,
			expectCnt: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSitus, _, _, _, svc := setupSitusService()
			tt.setupMock(mockSitus)

			res, err := svc.GetAllSitusForPublic(ctx, filter)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Len(t, res, tt.expectCnt)
			}

			mockSitus.AssertExpectations(t)
		})
	}
}

func TestSitusKeagamaanServiceImpl_GetSitusDetailForPublic(t *testing.T) {
	ctx := context.Background()
	situsId := uuid.New()

	tests := []struct {
		name      string
		setupMock func(*MockSitusKeagamaanRepo)
		expectErr bool
	}{
		{
			name: "Success - Get Detail Public",
			setupMock: func(ms *MockSitusKeagamaanRepo) {
				dummyDetail := &dto.SitusPublicDetailResponse{
					FasilitasJSON: []byte(`{"parkir": true, "toilet": true}`),
					GaleriJSON:    []byte(`["img1.jpg", "img2.jpg"]`),
				}
				ms.On("ReadDetailForPublic", ctx, situsId).Return(dummyDetail, nil).Once()
			},
			expectErr: false,
		},
		{
			name: "Error - Not Found",
			setupMock: func(ms *MockSitusKeagamaanRepo) {
				ms.On("ReadDetailForPublic", ctx, situsId).Return((*dto.SitusPublicDetailResponse)(nil), sql.ErrNoRows).Once()
			},
			expectErr: true,
		},
		{
			name: "Error - Repo Failed",
			setupMock: func(ms *MockSitusKeagamaanRepo) {
				ms.On("ReadDetailForPublic", ctx, situsId).Return((*dto.SitusPublicDetailResponse)(nil), errors.New("db error")).Once()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSitus, _, _, _, svc := setupSitusService()
			tt.setupMock(mockSitus)

			res, err := svc.GetSitusDetailForPublic(ctx, situsId)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				// Verify JSON unmarshaling worked
				assert.NotNil(t, res.Fasilitas)
				assert.NotNil(t, res.Galeri)
			}

			mockSitus.AssertExpectations(t)
		})
	}
}

func TestSitusKeagamaanServiceImpl_GetLandingStats(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		setupMock func(*MockSitusKeagamaanRepo)
		expectErr bool
	}{
		{
			name: "Success - Get Landing Stats",
			setupMock: func(ms *MockSitusKeagamaanRepo) {
				ms.On("GetLandingStats", ctx).Return(&dto.LandingStatsResponse{}, nil).Once()
			},
			expectErr: false,
		},
		{
			name: "Error - Repo Failed",
			setupMock: func(ms *MockSitusKeagamaanRepo) {
				ms.On("GetLandingStats", ctx).Return((*dto.LandingStatsResponse)(nil), errors.New("db error")).Once()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSitus, _, _, _, svc := setupSitusService()
			tt.setupMock(mockSitus)

			res, err := svc.GetLandingStats(ctx)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}

			mockSitus.AssertExpectations(t)
		})
	}
}
