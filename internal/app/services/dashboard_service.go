package services

import (
	"context"
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/repositories"
	"situs-keagamaan/internal/dto"
)

type DashboardService interface {
	GetActivities(ctx context.Context) ([]dto.ActivityResponse, error)
}

type dashboardServiceImpl struct {
	activityRepo repositories.ActivityRepository
}

func NewDashboardService(activityRepo repositories.ActivityRepository) DashboardService {
	return &dashboardServiceImpl{activityRepo: activityRepo}
}

func (s *dashboardServiceImpl) GetActivities(ctx context.Context) ([]dto.ActivityResponse, error) {
	activities, err := s.activityRepo.GetActivities(ctx)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}
	return activities, nil
}
