package services_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"sicemas/internal/app/services"
	"sicemas/internal/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestActivityServiceImpl_GetActivityByDate(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	ctx := context.Background()

	tests := []struct {
		name           string
		dateInput      string
		mockSetup      func(m *MockActivityRepo)
		expectedResult []dto.ActivityResponse
		expectError    bool
	}{
		{
			name:      "Success - Returns Data",
			dateInput: "2026-05-08",
			mockSetup: func(m *MockActivityRepo) {
				m.On("GetActivities", mock.Anything, mock.Anything, mock.Anything).
					Return([]dto.ActivityResponse{{}}, nil).Once()
			},
			expectedResult: []dto.ActivityResponse{{}},
			expectError:    false,
		},
		{
			name:      "Success - No Activities Found",
			dateInput: "2026-05-08",
			mockSetup: func(m *MockActivityRepo) {
				m.On("GetActivities", mock.Anything, mock.Anything, mock.Anything).
					Return([]dto.ActivityResponse(nil), nil).Once()
			},
			expectedResult: []dto.ActivityResponse{},
			expectError:    false,
		},
		{
			name:      "Error - Invalid Date Format",
			dateInput: "08-05-2026",
			mockSetup: func(m *MockActivityRepo) {
			},
			expectedResult: nil,
			expectError:    true,
		},
		{
			name:      "Error - Repository Returns Error",
			dateInput: "2026-05-08",
			mockSetup: func(m *MockActivityRepo) {
				m.On("GetActivities", mock.Anything, mock.Anything, mock.Anything).
					Return([]dto.ActivityResponse(nil), errors.New("database connection failed")).Once()
			},
			expectedResult: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockActivityRepo)

			tt.mockSetup(repo)

			service := services.NewActivityService(repo, logger)

			result, err := service.GetActivityByDate(ctx, tt.dateInput)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResult, result)

			repo.AssertExpectations(t)
		})
	}
}
