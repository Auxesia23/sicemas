package handlers

import (
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/services"
	"situs-keagamaan/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PolicyHandler interface {
	AddPolicy(c *fiber.Ctx) error
	RemovePolicy(c *fiber.Ctx) error
	GetFilteredPolicy(c *fiber.Ctx) error
}

type policyHandlerImpl struct {
	policyService services.PolicyService
	validate      *validator.Validate
}

func NewPolicyHandler(policyService services.PolicyService, validate *validator.Validate) PolicyHandler {
	return &policyHandlerImpl{
		policyService,
		validate,
	}
}

func (h *policyHandlerImpl) AddPolicy(c *fiber.Ctx) error {
	claim := c.Locals("claim").(*dto.AccessToken)
	actorId, err := uuid.Parse(claim.Subject)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	var body dto.PolicyRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.policyService.AddPolicy(body, actorId); err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (h *policyHandlerImpl) RemovePolicy(c *fiber.Ctx) error {
	claim := c.Locals("claim").(*dto.AccessToken)
	actorId, err := uuid.Parse(claim.Subject)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	var body dto.PolicyRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.policyService.RemovePolicy(body, actorId); err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *policyHandlerImpl) GetFilteredPolicy(c *fiber.Ctx) error {
	filter := c.Query("filter")

	policies, err := h.policyService.GetFilteredPolicy(filter)
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}

	return c.Status(200).JSON(policies)
}
