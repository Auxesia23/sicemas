package handlers

import (
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/services"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler interface {
	GetActivities(c *fiber.Ctx) error
}

type dashboardHandlerImpl struct {
	dashboardService services.DashboardService
}

func NewDashboardHandler(dashboardService services.DashboardService) *dashboardHandlerImpl {
	return &dashboardHandlerImpl{dashboardService: dashboardService}
}

func (h *dashboardHandlerImpl) GetActivities(c *fiber.Ctx) error {
	activities, err := h.dashboardService.GetActivities(c.Context())
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Message)
	}

	return c.JSON(activities)
}
