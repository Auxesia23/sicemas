package handlers

import (
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/services"
	"situs-keagamaan/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SitusKeagamaanHandler interface {
	CreateSitus(c *fiber.Ctx) error
	GetAllSitus(c *fiber.Ctx) error
}

type situsKeagamaanHandlerImpl struct {
	situsService services.SitusKeagamaanService
	validate     *validator.Validate
}

func NewSitusKeagamaanHandler(situsService services.SitusKeagamaanService, validate *validator.Validate) SitusKeagamaanHandler {
	return &situsKeagamaanHandlerImpl{situsService: situsService, validate: validate}
}

func (h *situsKeagamaanHandlerImpl) CreateSitus(c *fiber.Ctx) error {
	var body dto.SitusKeagamaanRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	user := c.Locals("claim").(*dto.AccessToken)
	userId, err := uuid.Parse(user.Subject)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	err = h.situsService.CreateSitusKeagamaan(c.Context(), &body, userId)
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Message)
	}

	return c.SendStatus(201)
}

func (h *situsKeagamaanHandlerImpl) GetAllSitus(c *fiber.Ctx) error {
	situs, err := h.situsService.GetAllSitusKeagamaan(c.Context())
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Message)
	}

	return c.Status(200).JSON(situs)
}
