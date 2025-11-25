package handlers

import (
	"situs-keagamaan/internal/models"
	"situs-keagamaan/internal/services"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	RegisterUser(c *fiber.Ctx) error
	LoginUser(c *fiber.Ctx) error
	VerifyOTP(c *fiber.Ctx) error
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

func (h *userHandlerImpl) RegisterUser(c *fiber.Ctx) error {
	var body models.UserRegister
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.userService.RegisterUser(c.Context(), &body); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *userHandlerImpl) LoginUser(c *fiber.Ctx) error {
	var body models.UserLogin
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.userService.LoginUser(c.Context(), &body); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.SendStatus(200)
}

func (h *userHandlerImpl) VerifyOTP(c *fiber.Ctx) error {
	var body models.UserVerifyOTP
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	token, err := h.userService.VerifyOTP(c.Context(), &body)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     "acces_token",
		Value:    token.AccessToken,
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute),
		HTTPOnly: true,
		Secure:   true,
	})
	return c.SendStatus(200)
}
