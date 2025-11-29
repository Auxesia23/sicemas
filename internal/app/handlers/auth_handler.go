package handlers

import (
	"situs-keagamaan/internal/app/services"
	"situs-keagamaan/internal/dto"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Login(c *fiber.Ctx) error
	VerifyOTP(c *fiber.Ctx) error
}

type authHandlerImpl struct {
	authService services.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService services.AuthService, validate *validator.Validate) AuthHandler {
	return &authHandlerImpl{
		authService,
		validate,
	}
}

func (h *authHandlerImpl) Login(c *fiber.Ctx) error {
	var body dto.UserLogin
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.authService.Login(c.Context(), &body); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.SendStatus(200)
}

func (h *authHandlerImpl) VerifyOTP(c *fiber.Ctx) error {
	var body dto.UserVerifyOTP
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	token, err := h.authService.VerifyOTP(c.Context(), &body)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     "acces_token",
		Value:    token.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute),
		HTTPOnly: true,
		Secure:   true,
	})
	return c.Status(200).JSON(fiber.Map{
		"access_token": token.AccessToken,
	})
}
