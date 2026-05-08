package services

import (
	"context"
	"log/slog"
	apperror "sicemas/internal/app/appError"
	"sicemas/internal/app/repositories"
	"sicemas/internal/dto"
	"time"
)

type ActivityService interface {
	GetActivityByDate(ctx context.Context, date string) ([]dto.ActivityResponse, error)
}

type ActivityServiceImpl struct {
	activityRepo repositories.ActivityRepository
	logger       *slog.Logger
}

func NewActivityService(activityRepo repositories.ActivityRepository, logger *slog.Logger) ActivityService {
	return &ActivityServiceImpl{activityRepo: activityRepo, logger: logger}
}

func (s *ActivityServiceImpl) GetActivityByDate(ctx context.Context, date string) ([]dto.ActivityResponse, error) {
	startDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		s.logger.Warn("Date parsing failed", "error", err)
		return nil, apperror.NewBadRequest("Format tanggal tidak valid, gunakan format YYYY-MM-DD")
	}
	endDate := startDate.AddDate(0, 0, 1)
	activities, err := s.activityRepo.GetActivities(ctx, startDate, endDate)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan")
	}
	if len(activities) == 0 {
		s.logger.Info("no activities found for date", "date", date)
		return []dto.ActivityResponse{}, nil
	}
	return activities, nil
}
