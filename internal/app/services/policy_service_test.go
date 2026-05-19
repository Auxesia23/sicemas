package services_test

import (
	"io"
	"log/slog"
	"testing"

	"sicemas/internal/app/services"
	"sicemas/internal/dto"

	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupPolicyService() (*MockActivityRepo, *casbin.Enforcer, services.PolicyService) {
	mockActivityRepo := new(MockActivityRepo)
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

	svc := services.NewPolicyService(
		enforcer,
		mockActivityRepo,
		logger,
	)

	return mockActivityRepo, enforcer, svc
}

func TestPolicyServiceImpl_AddPolicy(t *testing.T) {
	actorId := uuid.New()
	req := dto.PolicyRequest{
		Subject: "admin",
		Object:  "users",
		Action:  "read",
	}

	tests := []struct {
		name          string
		mockSetup     func(mockAct *MockActivityRepo, enforcer *casbin.Enforcer)
		expectedError bool
		errorMsg      string
	}{
		{
			name: "Success - Add New Policy",
			mockSetup: func(mockAct *MockActivityRepo, enforcer *casbin.Enforcer) {
				mockAct.On("InsertActivity", mock.Anything, mock.AnythingOfType("*entity.Activity")).
					Return(nil).Once()
			},
			expectedError: false,
		},
		{
			name: "Error - Policy Already Exists",
			mockSetup: func(mockAct *MockActivityRepo, enforcer *casbin.Enforcer) {
				enforcer.AddPolicy(req.Subject, req.Object, req.Action)
			},
			expectedError: true,
			errorMsg:      "Policy sudah ada!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAct, enforcer, svc := setupPolicyService()

			tt.mockSetup(mockAct, enforcer)

			err := svc.AddPolicy(req, actorId)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			mockAct.AssertExpectations(t)
		})
	}
}

func TestPolicyServiceImpl_GetFilteredPolicy(t *testing.T) {
	filter := "admin"

	tests := []struct {
		name          string
		mockSetup     func(enforcer *casbin.Enforcer)
		expectedCount int
		expectedError bool
	}{
		{
			name: "Success - Returns Filtered Policies",
			mockSetup: func(enforcer *casbin.Enforcer) {
				enforcer.AddPolicy("admin", "users", "read")
				enforcer.AddPolicy("admin", "users", "write")
				enforcer.AddPolicy("user", "users", "read")
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name: "Success - Empty Result",
			mockSetup: func(enforcer *casbin.Enforcer) {
			},
			expectedCount: 0,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, enforcer, svc := setupPolicyService()

			tt.mockSetup(enforcer)

			result, err := svc.GetFilteredPolicy(filter)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedCount, len(result))
			}
		})
	}
}

func TestPolicyServiceImpl_RemovePolicy(t *testing.T) {
	actorId := uuid.New()
	req := dto.PolicyRequest{
		Subject: "admin",
		Object:  "users",
		Action:  "read",
	}

	tests := []struct {
		name          string
		mockSetup     func(mockAct *MockActivityRepo, enforcer *casbin.Enforcer)
		expectedError bool
		errorMsg      string
	}{
		{
			name: "Success - Remove Existing Policy",
			mockSetup: func(mockAct *MockActivityRepo, enforcer *casbin.Enforcer) {
				enforcer.AddPolicy(req.Subject, req.Object, req.Action)

				mockAct.On("InsertActivity", mock.Anything, mock.AnythingOfType("*entity.Activity")).
					Return(nil).Once()
			},
			expectedError: false,
		},
		{
			name: "Error - Policy Not Found",
			mockSetup: func(mockAct *MockActivityRepo, enforcer *casbin.Enforcer) {
			},
			expectedError: true,
			errorMsg:      "Policy tidak ditemukan!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAct, enforcer, svc := setupPolicyService()

			tt.mockSetup(mockAct, enforcer)

			err := svc.RemovePolicy(req, actorId)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			mockAct.AssertExpectations(t)
		})
	}
}
