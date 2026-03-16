package services

import (
	"context"
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/repositories"
	"situs-keagamaan/internal/dto"
)

type DashboardService interface {
	GetDashboardData(ctx context.Context) (*dto.DashboardResponse, error)
}

type dashboardServiceImpl struct {
	activityRepo  repositories.ActivityRepository
	dashboardRepo repositories.DashboardRepository
}

func NewDashboardService(activityRepo repositories.ActivityRepository, dashboardRepo repositories.DashboardRepository) DashboardService {
	return &dashboardServiceImpl{activityRepo: activityRepo, dashboardRepo: dashboardRepo}
}

func (s *dashboardServiceImpl) GetDashboardData(ctx context.Context) (*dto.DashboardResponse, error) {
	activities, err := s.activityRepo.GetActivities(ctx)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}

	stats, err := s.dashboardRepo.GetDashboardStats(ctx)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}
	recentSites, err := s.dashboardRepo.GetRecentSites(ctx)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}
	statistikJenis, err := s.dashboardRepo.GetStatistikJenis(ctx)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}

	response := &dto.DashboardResponse{
		Stats:            stats,
		StatistikJenis:   statistikJenis,
		RecentActivities: activities,
		RecentSites:      recentSites,
	}

	return response, nil
}
