package handlers

import (
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/services"
	"situs-keagamaan/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type JenisSitusHandler interface {
	CreateJenisSitus(c *fiber.Ctx) error
	UpdateJenisSitus(c *fiber.Ctx) error
	DeleteJenisSitus(c *fiber.Ctx) error
	GetAllJenisSitus(c *fiber.Ctx) error
}

type jenisSitusHandlerImpl struct {
	service  services.JenisSitusService
	validate *validator.Validate
}

func NewJenisSitusHandler(service services.JenisSitusService, validate *validator.Validate) JenisSitusHandler {
	return &jenisSitusHandlerImpl{service: service, validate: validate}
}

func (h *jenisSitusHandlerImpl) CreateJenisSitus(c *fiber.Ctx) error {
	var body dto.JenisSitusRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	newJenisSitus, err := h.service.CreateJenisSitus(c.Context(), &body)
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}
	return c.JSON(newJenisSitus)
}

func (h *jenisSitusHandlerImpl) GetAllJenisSitus(c *fiber.Ctx) error {
	jenisSitus, err := h.service.GetAllJenisSitus(c.Context())
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}
	return c.JSON(jenisSitus)
}

func (h *jenisSitusHandlerImpl) UpdateJenisSitus(c *fiber.Ctx) error {
	var body dto.JenisSitusRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	id := c.Params("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	err = h.service.UpdateJenisSitus(c.Context(), uuid, &body)
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}
	return c.SendStatus(200)
}

func (h *jenisSitusHandlerImpl) DeleteJenisSitus(c *fiber.Ctx) error {
	id := c.Params("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	err = h.service.DeleteJenisSitus(c.Context(), uuid)
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}
	return c.SendStatus(200)
}
