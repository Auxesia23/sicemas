package handlers

import (
	"situs-keagamaan/internal/app/services"
	"situs-keagamaan/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	Register(c *fiber.Ctx) error
}

type userHandlerImpl struct {
	userService services.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService services.UserService, validate *validator.Validate) UserHandler {
	return &userHandlerImpl{
		userService,
		validate,
	}
}

func (h *userHandlerImpl) Register(c *fiber.Ctx) error {
	var body dto.UserRegister
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.userService.Register(c.Context(), &body); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}
