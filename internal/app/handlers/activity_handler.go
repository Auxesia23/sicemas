package handlers

import (
	"log"
	apperror "sicemas/internal/app/appError"
	"sicemas/internal/app/services"

	"github.com/gofiber/fiber/v2"
)

type ActivityHandler interface {
	GetActivityByDate(c *fiber.Ctx) error
}

type ActivityHandlerImpl struct {
	activityService services.ActivityService
}

func NewActivityHandler(activityService services.ActivityService) ActivityHandler {
	return &ActivityHandlerImpl{activityService: activityService}
}

func (h *ActivityHandlerImpl) GetActivityByDate(c *fiber.Ctx) error {
	date := c.Query("date")
	log.Printf("date: %s", date)
	activities, err := h.activityService.GetActivityByDate(c.Context(), date)
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}
	return c.JSON(activities)
}
