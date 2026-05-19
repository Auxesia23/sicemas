package handlers

import (
	apperror "sicemas/internal/app/appError"
	"sicemas/internal/app/services"
	"sicemas/internal/dto"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler interface {
	Login(c *fiber.Ctx) error
	VerifyOTP(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	VerifyStepUpOTP(c *fiber.Ctx) error
	ResendStepUpOTP(c *fiber.Ctx) error
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
		SameSite: "Strict",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token.AccessToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * 15),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "csrf_token",
		Value:    token.CSRFToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * 15),
		HTTPOnly: false,
		Secure:   true,
		SameSite: "Strict",
	})
	return c.SendStatus(200)
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
		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Strict",
		})
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    "",
			Path:     "/",
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Strict",
		})
		c.Cookie(&fiber.Cookie{
			Name:     "csrf_token",
			Value:    "",
			Path:     "/",
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: false,
			Secure:   true,
			SameSite: "Strict",
		})
		return c.Status(e.Status).SendString(e.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    token.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token.AccessToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * 15),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "csrf_token",
		Value:    token.CSRFToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * 15),
		HTTPOnly: false,
		Secure:   true,
		SameSite: "Strict",
	})
	return c.SendStatus(200)
}

func (h *authHandlerImpl) Logout(c *fiber.Ctx) error {
	accessToken := c.Cookies("access_token")
	var accessTokenPtr *string
	if accessToken == "" {
		accessTokenPtr = nil
	} else {
		accessTokenPtr = &accessToken
	}

	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(200).SendString("Kamu belum login si..")
	}

	if err := h.authService.Logout(c.Context(), refreshToken, accessTokenPtr); err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: false,
		Secure:   true,
		SameSite: "Strict",
	})
	return c.SendStatus(200)
}

func (h *authHandlerImpl) VerifyStepUpOTP(c *fiber.Ctx) error {
	requestContext := c.Locals("context").(*dto.SessionRequest)

	var body dto.UserStepUpOTP
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	user := c.Locals("claim").(*dto.AccessToken)
	SID, err := uuid.Parse(user.SID)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if err := h.authService.VerifyStepUpOTP(c.Context(), user.Subject, SID, body.OTP, requestContext); err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}

	return c.SendStatus(200)
}

func (h *authHandlerImpl) ResendStepUpOTP(c *fiber.Ctx) error {
	user := c.Locals("claim").(*dto.AccessToken)
	if err := h.authService.TriggerStepUpOTP(c.Context(), user.Subject); err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Error())
	}
	return c.SendStatus(200)
}
