package handlers

import (
	apperror "sicemas/internal/app/appError"
	"sicemas/internal/app/services"
	"sicemas/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RoleHandler interface {
	CreateRole(c *fiber.Ctx) error
	GetAllRole(c *fiber.Ctx) error
	DeleteRole(c *fiber.Ctx) error
}

type roleHandlerImpl struct {
	roleService services.RoleService
	validate    *validator.Validate
}

func NewRoleHandler(roleService services.RoleService, validate *validator.Validate) RoleHandler {
	return &roleHandlerImpl{
		roleService,
		validate,
	}
}

func (h *roleHandlerImpl) CreateRole(c *fiber.Ctx) error {
	claim := c.Locals("claim").(*dto.AccessToken)
	actorId, err := uuid.Parse(claim.Subject)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	var body dto.RoleRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	newRole, err := h.roleService.CreateRole(c.Context(), body.Name, actorId)
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}

	return c.Status(201).JSON(newRole)
}

func (h *roleHandlerImpl) GetAllRole(c *fiber.Ctx) error {
	roles, err := h.roleService.GetAllRole(c.Context())
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}

	return c.Status(200).JSON(roles)
}

func (h *roleHandlerImpl) DeleteRole(c *fiber.Ctx) error {
	claim := c.Locals("claim").(*dto.AccessToken)
	actorId, err := uuid.Parse(claim.Subject)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	id := c.Params("roleId")
	if err := h.roleService.DeleteRole(c.Context(), id, actorId); err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
