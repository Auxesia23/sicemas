package services_test

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"sicemas/internal/app/services"
	"sicemas/internal/entity"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRoleService() (*MockRoleRepo, *MockActivityRepo, services.RoleService) {
	mockRoleRepo := new(MockRoleRepo)
	mockActivityRepo := new(MockActivityRepo)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	svc := services.NewRoleService(
		mockRoleRepo,
		mockActivityRepo,
		logger,
	)

	return mockRoleRepo, mockActivityRepo, svc
}

func TestRoleServiceImpl_CreateRole(t *testing.T) {
	defer func() {
		os.Unsetenv("AES_256_KEY")
		os.Unsetenv("PEPPER")
	}()
	os.Setenv("AES_256_KEY", "MTIzNDU2Nzg5MDEyMzQ1Njc4OTA=")
	os.Setenv("PEPPER", "dummy-pepper")

	ctx := context.Background()
	actorId := uuid.New()

	tests := []struct {
		name       string
		roleName   string
		setupMock  func(*MockRoleRepo, *MockActivityRepo)
		expectErr  bool
		expectName string
	}{
		{
			name:     "Success - Create Role",
			roleName: "ADMIN",
			setupMock: func(mr *MockRoleRepo, ma *MockActivityRepo) {
				newRole := &entity.Role{
					ID:        uuid.New().String(),
					Name:      "admin",
					CreatedAt: time.Now(),
				}
				mr.On("Create", ctx, "admin").Return(newRole, nil).Once()
				ma.On("InsertActivity", ctx, mock.AnythingOfType("*entity.Activity")).Return(nil).Once()
			},
			expectErr:  false,
			expectName: "admin",
		},
		{
			name:     "Error - Role Already Exists",
			roleName: "ADMIN",
			setupMock: func(mr *MockRoleRepo, ma *MockActivityRepo) {
				mr.On("Create", ctx, "admin").
					Return((*entity.Role)(nil), &pq.Error{Code: "23505"}).Once()
			},
			expectErr:  true,
			expectName: "",
		},
		{
			name:     "Error - Create Failed",
			roleName: "ADMIN",
			setupMock: func(mr *MockRoleRepo, ma *MockActivityRepo) {
				mr.On("Create", ctx, "admin").
					Return((*entity.Role)(nil), errors.New("database error")).Once()
			},
			expectErr:  true,
			expectName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoleRepo, mockActivityRepo, svc := setupRoleService()
			mockRoleRepo.ExpectedCalls = nil
			mockActivityRepo.ExpectedCalls = nil

			tt.setupMock(mockRoleRepo, mockActivityRepo)

			result, err := svc.CreateRole(ctx, tt.roleName, actorId)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectName, result.Name)
			}

			mockRoleRepo.AssertExpectations(t)
			mockActivityRepo.AssertExpectations(t)
		})
	}
}

func TestRoleServiceImpl_GetAllRole(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		setupMock func(*MockRoleRepo)
		expectErr bool
		expectCnt int
	}{
		{
			name: "Success - Get All Roles",
			setupMock: func(mr *MockRoleRepo) {
				roles := []entity.Role{
					{ID: uuid.New().String(), Name: "admin", CreatedAt: time.Now()},
					{ID: uuid.New().String(), Name: "user", CreatedAt: time.Now()},
				}
				mr.On("ReadAll", ctx).Return(roles, nil).Once()
			},
			expectErr: false,
			expectCnt: 2,
		},
		{
			name: "Success - No Roles Found",
			setupMock: func(mr *MockRoleRepo) {
				mr.On("ReadAll", ctx).Return([]entity.Role{}, nil).Once()
			},
			expectErr: false,
			expectCnt: 0,
		},
		{
			name: "Error - Repository Failed",
			setupMock: func(mr *MockRoleRepo) {
				mr.On("ReadAll", ctx).
					Return([]entity.Role(nil), errors.New("database error")).Once()
			},
			expectErr: true,
			expectCnt: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoleRepo, _, svc := setupRoleService()
			mockRoleRepo.ExpectedCalls = nil

			tt.setupMock(mockRoleRepo)

			result, err := svc.GetAllRole(ctx)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectCnt, len(result))
			}

			mockRoleRepo.AssertExpectations(t)
		})
	}
}

func TestRoleServiceImpl_DeleteRole(t *testing.T) {
	ctx := context.Background()
	actorId := uuid.New()
	roleId := uuid.New().String()

	tests := []struct {
		name      string
		setupMock func(*MockRoleRepo, *MockActivityRepo)
		expectErr bool
	}{
		{
			name: "Success - Delete Role",
			setupMock: func(mr *MockRoleRepo, ma *MockActivityRepo) {
				role := &entity.Role{
					ID:        roleId,
					Name:      "admin",
					CreatedAt: time.Now(),
				}
				mr.On("ReadOne", ctx, roleId).Return(role, nil).Once()
				mr.On("Delete", ctx, roleId).Return(nil).Once()
				ma.On("InsertActivity", ctx, mock.AnythingOfType("*entity.Activity")).Return(nil).Once()
			},
			expectErr: false,
		},
		{
			name: "Error - Role Not Found When Reading",
			setupMock: func(mr *MockRoleRepo, ma *MockActivityRepo) {
				mr.On("ReadOne", ctx, roleId).Return((*entity.Role)(nil), sql.ErrNoRows).Once()
			},
			expectErr: true,
		},
		{
			name: "Error - Database Error When Reading",
			setupMock: func(mr *MockRoleRepo, ma *MockActivityRepo) {
				mr.On("ReadOne", ctx, roleId).Return((*entity.Role)(nil), errors.New("database error")).Once()
			},
			expectErr: true,
		},
		{
			name: "Error - Delete Failed",
			setupMock: func(mr *MockRoleRepo, ma *MockActivityRepo) {
				role := &entity.Role{
					ID:        roleId,
					Name:      "admin",
					CreatedAt: time.Now(),
				}
				mr.On("ReadOne", ctx, roleId).Return(role, nil).Once()
				mr.On("Delete", ctx, roleId).Return(errors.New("database error")).Once()
			},
			expectErr: true,
		},
		{
			name: "Error - Role Not Found During Delete",
			setupMock: func(mr *MockRoleRepo, ma *MockActivityRepo) {
				role := &entity.Role{
					ID:        roleId,
					Name:      "admin",
					CreatedAt: time.Now(),
				}
				mr.On("ReadOne", ctx, roleId).Return(role, nil).Once()
				mr.On("Delete", ctx, roleId).Return(sql.ErrNoRows).Once()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoleRepo, mockActivityRepo, svc := setupRoleService()
			mockRoleRepo.ExpectedCalls = nil
			mockActivityRepo.ExpectedCalls = nil

			tt.setupMock(mockRoleRepo, mockActivityRepo)

			err := svc.DeleteRole(ctx, roleId, actorId)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRoleRepo.AssertExpectations(t)
			mockActivityRepo.AssertExpectations(t)
		})
	}
}
