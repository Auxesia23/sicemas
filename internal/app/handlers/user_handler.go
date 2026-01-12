package handlers

import (
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/services"
	"situs-keagamaan/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler interface {
	Register(c *fiber.Ctx) error
	GetAllUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
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
	var body dto.UserRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.userService.Register(c.Context(), &body); err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *userHandlerImpl) GetAllUser(c *fiber.Ctx) error {
	users, err := h.userService.GetAllUser(c.Context())
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}
	return c.Status(200).JSON(users)
}

func (h *userHandlerImpl) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("userId")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	var body dto.UserRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.userService.UpdateUser(c.Context(), uuid, &body); err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}

	return c.SendStatus(200)
}

func (h *userHandlerImpl) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("userId")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.userService.DeleteUser(c.Context(), uuid); err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Message)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
