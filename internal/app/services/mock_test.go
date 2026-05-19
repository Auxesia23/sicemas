package services_test

import (
	"context"
	"time"

	"sicemas/internal/dto"
	"sicemas/internal/entity"

	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// ==========================================
// MOCK USER REPO
// ==========================================

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, in *dto.UserRequest, index []byte) (string, error) {
	args := m.Called(ctx, in, index)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepo) ReadAll(ctx context.Context) ([]entity.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.User), args.Error(1)
}

func (m *MockUserRepo) ReadOne(ctx context.Context, index []byte) (*entity.User, error) {
	args := m.Called(ctx, index)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepo) ReadById(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, id uuid.UUID, in *dto.UserRequest, newIndex []byte) error {
	args := m.Called(ctx, id, in, newIndex)
	return args.Error(0)
}

func (m *MockUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ==========================================
// MOCK ACTIVITY REPO
// ==========================================

type MockActivityRepo struct {
	mock.Mock
}

func (m *MockActivityRepo) InsertActivity(ctx context.Context, act *entity.Activity) error {
	args := m.Called(ctx, act)
	return args.Error(0)
}

func (m *MockActivityRepo) GetActivities(ctx context.Context, startDate, endDate time.Time) ([]dto.ActivityResponse, error) {
	args := m.Called(ctx, startDate, endDate)
	if args.Get(0) != nil {
		return args.Get(0).([]dto.ActivityResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

// ==========================================
// MOCK ROLE REPO
// ==========================================

type MockRoleRepo struct {
	mock.Mock
}

func (m *MockRoleRepo) Create(ctx context.Context, name string) (*entity.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Role), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRoleRepo) ReadAll(ctx context.Context) ([]entity.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]entity.Role), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRoleRepo) ReadOne(ctx context.Context, id string) (*entity.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Role), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRoleRepo) Delete(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

// ==========================================
// MOCK JENIS SITUS REPO
// ==========================================

type MockJenisSitusRepo struct {
	mock.Mock
}

func (m *MockJenisSitusRepo) Create(ctx context.Context, in *dto.JenisSitusRequest) (*entity.JenisSitus, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.JenisSitus), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockJenisSitusRepo) ReadAll(ctx context.Context) ([]dto.JenisSitusResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]dto.JenisSitusResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockJenisSitusRepo) ReadOne(ctx context.Context, id uuid.UUID) (*dto.JenisSitusResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.JenisSitusResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockJenisSitusRepo) Update(ctx context.Context, id uuid.UUID, in *dto.JenisSitusRequest) error {
	args := m.Called(ctx, id, in)
	return args.Error(0)
}

func (m *MockJenisSitusRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ==========================================
// MOCK DASHBOARD REPO
// ==========================================

type MockDashboardRepo struct {
	mock.Mock
}

func (m *MockDashboardRepo) GetDashboardStats(ctx context.Context) (*dto.Stats, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.Stats), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDashboardRepo) GetRecentSites(ctx context.Context) ([]dto.SitusKeagamaanResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]dto.SitusKeagamaanResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDashboardRepo) GetStatistikJenis(ctx context.Context) ([]dto.StatistikJenis, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]dto.StatistikJenis), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDashboardRepo) GetRecentActivities(ctx context.Context) ([]dto.ActivityResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]dto.ActivityResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

// ==========================================
// MOCK CACHE
// ==========================================

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Set(ctx context.Context, key string, data any, ttl time.Duration) error {
	args := m.Called(ctx, key, data, ttl)
	return args.Error(0)
}

func (m *MockCache) Get(ctx context.Context, key string, dest any) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

// ==========================================
// DUMMY CASBIN ADAPTER
// ==========================================

type DummyCasbinAdapter struct{}

func (d *DummyCasbinAdapter) LoadPolicy(model casbinModel.Model) error                   { return nil }
func (d *DummyCasbinAdapter) SavePolicy(model casbinModel.Model) error                   { return nil }
func (d *DummyCasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error    { return nil }
func (d *DummyCasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error { return nil }
func (d *DummyCasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}

// ==========================================
// MOCK ENFORCER
// ==========================================

type MockEnforcer struct {
	mock.Mock
}

func (m *MockEnforcer) AddPolicy(params ...string) (bool, error) {
	args := m.Called(params)
	return args.Bool(0), args.Error(1)
}

func (m *MockEnforcer) RemovePolicy(params ...string) (bool, error) {
	args := m.Called(params)
	return args.Bool(0), args.Error(1)
}

func (m *MockEnforcer) GetFilteredPolicy(fieldIndex int, fieldValue string) ([][]string, error) {
	args := m.Called(fieldIndex, fieldValue)
	if args.Get(0) != nil {
		return args.Get(0).([][]string), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockEnforcer) LoadPolicy() error {
	return nil
}

// ==========================================
// MOCK: SitusKeagamaanRepository
// ==========================================

type MockSitusKeagamaanRepo struct {
	mock.Mock
}

func (m *MockSitusKeagamaanRepo) Create(ctx context.Context, in *dto.SitusKeagamaanRequest, author uuid.UUID) (uuid.UUID, error) {
	args := m.Called(ctx, in, author)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockSitusKeagamaanRepo) ReadAll(ctx context.Context) ([]dto.SitusKeagamaanResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]dto.SitusKeagamaanResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSitusKeagamaanRepo) ReadOwn(ctx context.Context, userId uuid.UUID) ([]dto.SitusKeagamaanResponse, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) != nil {
		return args.Get(0).([]dto.SitusKeagamaanResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSitusKeagamaanRepo) ReadDetail(ctx context.Context, id uuid.UUID) (*dto.SitusKeagamaanDetailResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.SitusKeagamaanDetailResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSitusKeagamaanRepo) ReadAllDetail(ctx context.Context, jenisSitus string) ([]dto.SitusKeagamaanDetailResponse, error) {
	args := m.Called(ctx, jenisSitus)
	if args.Get(0) != nil {
		return args.Get(0).([]dto.SitusKeagamaanDetailResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSitusKeagamaanRepo) Update(ctx context.Context, id uuid.UUID, in *dto.SitusKeagamaanUpdate) error {
	args := m.Called(ctx, id, in)
	return args.Error(0)
}

func (m *MockSitusKeagamaanRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSitusKeagamaanRepo) Verify(ctx context.Context, id uuid.UUID, in *dto.VerifikasiSitusRequest) error {
	args := m.Called(ctx, id, in)
	return args.Error(0)
}

func (m *MockSitusKeagamaanRepo) CheckOwnership(ctx context.Context, id uuid.UUID, userId uuid.UUID) error {
	args := m.Called(ctx, id, userId)
	return args.Error(0)
}

func (m *MockSitusKeagamaanRepo) ReadForPublic(ctx context.Context, filter dto.PublicListFilter) ([]dto.SitusPublicListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) != nil {
		return args.Get(0).([]dto.SitusPublicListResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSitusKeagamaanRepo) ReadDetailForPublic(ctx context.Context, id uuid.UUID) (*dto.SitusPublicDetailResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.SitusPublicDetailResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSitusKeagamaanRepo) GetLandingStats(ctx context.Context) (*dto.LandingStatsResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.LandingStatsResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

// ==========================================
// MOCK: FotoSitusRepository
// ==========================================

type MockFotoSitusRepo struct {
	mock.Mock
}

func (m *MockFotoSitusRepo) BulkCreate(ctx context.Context, data []dto.FotoSitusPayload) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockFotoSitusRepo) GetPublicIDs(ctx context.Context, id []uuid.UUID) ([]string, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).([]string), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFotoSitusRepo) GetBySitusID(ctx context.Context, situsId uuid.UUID) (*[]dto.FotoResponse, error) {
	args := m.Called(ctx, situsId)
	if args.Get(0) != nil {
		return args.Get(0).(*[]dto.FotoResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFotoSitusRepo) Delete(ctx context.Context, id []uuid.UUID, situsId uuid.UUID) error {
	args := m.Called(ctx, id, situsId)
	return args.Error(0)
}
