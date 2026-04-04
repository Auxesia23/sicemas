package handlers

import (
	apperror "sicemas/internal/app/appError"
	"sicemas/internal/app/services"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler interface {
	GetDashboardData(c *fiber.Ctx) error
}

type dashboardHandlerImpl struct {
	dashboardService services.DashboardService
}

func NewDashboardHandler(dashboardService services.DashboardService) *dashboardHandlerImpl {
	return &dashboardHandlerImpl{dashboardService: dashboardService}
}

func (h *dashboardHandlerImpl) GetDashboardData(c *fiber.Ctx) error {
	activities, err := h.dashboardService.GetDashboardData(c.Context())
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Message)
	}

	return c.JSON(activities)
}
