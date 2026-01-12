package handlers

import (
	"log"
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/services"
	"situs-keagamaan/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PolicyHandler interface {
	AddPolicy(c *fiber.Ctx) error
	RemovePolicy(c *fiber.Ctx) error
	GetFilteredPolicy(c *fiber.Ctx) error
	AddGroupPolicy(c *fiber.Ctx) error
	RemoveGroupPolicy(c *fiber.Ctx) error
	GetFilteredGroupPolicy(c *fiber.Ctx) error
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
	var body dto.PolicyRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.policyService.AddPolicy(body); err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (h *policyHandlerImpl) RemovePolicy(c *fiber.Ctx) error {
	var body dto.PolicyRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.policyService.RemovePolicy(body); err != nil {
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

func (h *policyHandlerImpl) AddGroupPolicy(c *fiber.Ctx) error {
	var body dto.PolicyGroupRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.policyService.AddGroupPolicy(body); err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (h *policyHandlerImpl) RemoveGroupPolicy(c *fiber.Ctx) error {
	var body dto.PolicyGroupRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.policyService.RemoveGroupPolicy(body); err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *policyHandlerImpl) GetFilteredGroupPolicy(c *fiber.Ctx) error {
	filter := c.Query("filter")
	log.Println(filter)
	policies, err := h.policyService.GetFilteredGroupPolicy(filter)
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}

	return c.Status(200).JSON(policies)
}
