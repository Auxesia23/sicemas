package services

import (
	"context"
	"log/slog"
	apperror "sicemas/internal/app/appError"
	"sicemas/internal/app/repositories"
	"sicemas/internal/dto"
)

type DashboardService interface {
	GetDashboardData(ctx context.Context) (*dto.DashboardResponse, error)
}

type dashboardServiceImpl struct {
	activityRepo  repositories.ActivityRepository
	dashboardRepo repositories.DashboardRepository
	logger        *slog.Logger
}

func NewDashboardService(
	activityRepo repositories.ActivityRepository,
	dashboardRepo repositories.DashboardRepository,
	logger *slog.Logger,
) DashboardService {
	return &dashboardServiceImpl{
		activityRepo:  activityRepo,
		dashboardRepo: dashboardRepo,
		logger:        logger,
	}
}

func (s *dashboardServiceImpl) GetDashboardData(ctx context.Context) (*dto.DashboardResponse, error) {
	s.logger.Info("fetching dashboard data")

	recentActivities, err := s.dashboardRepo.GetRecentActivities(ctx)
	if err != nil {
		s.logger.Error("failed to get activities for dashboard", "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}
	if recentActivities == nil {
		recentActivities = make([]dto.ActivityResponse, 0)
	}

	stats, err := s.dashboardRepo.GetDashboardStats(ctx)
	if err != nil {
		s.logger.Error("failed to get dashboard stats", "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}

	recentSites, err := s.dashboardRepo.GetRecentSites(ctx)
	if err != nil {
		s.logger.Error("failed to get recent sites", "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}
	if recentSites == nil {
		recentSites = make([]dto.SitusKeagamaanResponse, 0)
	}

	statistikJenis, err := s.dashboardRepo.GetStatistikJenis(ctx)
	if err != nil {
		s.logger.Error("failed to get statistik jenis", "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}
	if statistikJenis == nil {
		statistikJenis = make([]dto.StatistikJenis, 0)
	}

	response := &dto.DashboardResponse{
		Stats:            stats,
		StatistikJenis:   statistikJenis,
		RecentActivities: recentActivities,
		RecentSites:      recentSites,
	}

	s.logger.Info("dashboard data fetched successfully")
	return response, nil
}
