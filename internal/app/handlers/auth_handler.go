package handlers

import (
	"log"
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/services"
	"situs-keagamaan/internal/dto"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Login(c *fiber.Ctx) error
	VerifyOTP(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
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
	refreshToken := c.Cookies("refresh_token")
	if refreshToken != "" {
		return c.Status(fiber.StatusBadRequest).SendString("Anda sudah login dan dalam sesi aktif")
	}
	var body dto.UserLogin
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.authService.Login(c.Context(), &body); err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}

	return c.SendStatus(200)
}

func (h *authHandlerImpl) VerifyOTP(c *fiber.Ctx) error {
	context := c.Locals("context").(*dto.SessionRequest)
	var body dto.UserVerifyOTP
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	token, err := h.authService.VerifyOTP(c.Context(), &body, context)
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    token.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
		Secure:   true,
	})
	return c.Status(200).JSON(fiber.Map{
		"access_token": token.AccessToken,
	})
}

func (h *authHandlerImpl) Refresh(c *fiber.Ctx) error {
	requestContext := c.Locals("context").(*dto.SessionRequest)
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Refresh token tidak ada atau kadaluarsa. Silahkan login kembali!")
	}
	token, err := h.authService.RefreshToken(c.Context(), refreshToken, requestContext)
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    token.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
		Secure:   true,
	})
	return c.Status(200).JSON(fiber.Map{
		"access_token": token.AccessToken,
	})
}

func (h *authHandlerImpl) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	bearerToken := strings.Split(authHeader, " ")
	var token *string
	if authHeader == "" {
		token = nil
		log.Println("Token Kosong")
	} else {
		token = &bearerToken[1]
		log.Println("Token Ada")
	}

	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(200).SendString("Kamu belum login si..")
	}

	if err := h.authService.Logout(c.Context(), refreshToken, token); err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}
	c.Cookie(&fiber.Cookie{
		Name:    "refresh_token",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-time.Hour),
	})
	return c.SendStatus(200)
}
