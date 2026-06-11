package middlewares

import (
	"context"
	"net/netip"
	"time"

	"sicemas/internal/dto"
	"sicemas/internal/entity"
	"sicemas/internal/geoip"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

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
// MOCK LOCATOR
// ==========================================

type MockLocator struct {
	mock.Mock
}

func (m *MockLocator) Lookup(ip netip.Addr) (*geoip.Location, error) {
	args := m.Called(ip)
	if args.Get(0) != nil {
		return args.Get(0).(*geoip.Location), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLocator) Close() error {
	args := m.Called()
	return args.Error(0)
}

// ==========================================
// MOCK AUTH SERVICE
// ==========================================

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, in *dto.UserLogin) error {
	args := m.Called(ctx, in)
	return args.Error(0)
}

func (m *MockAuthService) VerifyOTP(ctx context.Context, in *dto.UserVerifyOTP, loginContext *dto.SessionRequest) (*dto.Token, error) {
	args := m.Called(ctx, in, loginContext)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.Token), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string, requestContext *dto.SessionRequest) (*dto.Token, error) {
	args := m.Called(ctx, refreshToken, requestContext)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.Token), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, refreshToken string, accessToken *string) error {
	args := m.Called(ctx, refreshToken, accessToken)
	return args.Error(0)
}

func (m *MockAuthService) TriggerStepUpOTP(ctx context.Context, userId string) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func (m *MockAuthService) VerifyStepUpOTP(ctx context.Context, userId string, sid uuid.UUID, userOtp string, loginContext *dto.SessionRequest) error {
	args := m.Called(ctx, userId, sid, userOtp, loginContext)
	return args.Error(0)
}

// ==========================================
// DUMMY CASBIN ADAPTER
// ==========================================

type DummyCasbinAdapter struct{}

// Helper types for enforcer setup
type MockEnforcer struct {
	mock.Mock
}

// Helper to satisfy *casbin.Enforcer interface needs in tests
type dummyModel struct{}

// ==========================================
// MOCK ENFORCER (simplified for middleware test)
// ==========================================

type fakeEnforcer struct{}

// User entity helper for tests
func testUserEntity() *entity.User {
	return &entity.User{
		ID:          uuid.New(),
		NamaLengkap: "Test User",
		NIP:         "1234567890",
	}
}

func testSessionValue() *dto.SessionValue {
	return &dto.SessionValue{
		UserID:        uuid.New(),
		SID:           uuid.New(),
		UserAgent:     "test-agent",
		IPAddress:     "192.168.1.1",
		GeoLocation:   "Jakarta",
		DeviceID:      "device-123",
		IsMFAVerified: true,
	}
}
