package handlers

import (
	"mime/multipart"
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
	GetDetailSitus(c *fiber.Ctx) error
	DeleteSitus(c *fiber.Ctx) error
	UploadFotoSitus(c *fiber.Ctx) error
	DeleteFotoSitus(c *fiber.Ctx) error
	VerifySitus(c *fiber.Ctx) error
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
	id, err := h.situsService.CreateSitusKeagamaan(c.Context(), &body, userId)
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Message)
	}

	return c.Status(201).JSON(fiber.Map{"id": id})
}

func (h *situsKeagamaanHandlerImpl) UploadFotoSitus(c *fiber.Ctx) error {
	situsIdStr := c.Params("id")
	situsId, err := uuid.Parse(situsIdStr)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	fileHeaders := form.File["images"]
	if len(fileHeaders) == 0 {
		return c.Status(400).SendString("Tidak ada gambar yang diunggah")
	}

	var files []multipart.File

	defer func() {
		for _, f := range files {
			if f != nil {
				f.Close()
			}
		}
	}()

	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Gagal membuka file gambar"})
		}

		files = append(files, file)
	}

	err = h.situsService.UploadFotoSitus(c.Context(), files, situsId)
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Message)
	}

	return c.Status(201).JSON(fiber.Map{"message": "Foto berhasil diunggah"})
}

func (h *situsKeagamaanHandlerImpl) GetAllSitus(c *fiber.Ctx) error {
	claim := c.Locals("claim").(*dto.AccessToken)
	userId, err := uuid.Parse(claim.Subject)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}
	situs, err := h.situsService.GetAllSitusKeagamaan(c.Context(), userId)
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Message)
	}

	return c.Status(200).JSON(situs)
}

func (h *situsKeagamaanHandlerImpl) GetDetailSitus(c *fiber.Ctx) error {
	claim := c.Locals("claim").(*dto.AccessToken)
	userId, err := uuid.Parse(claim.Subject)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	situsId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}
	situs, err := h.situsService.GetDetailSitusKeagamaan(c.Context(), situsId, userId)
	if err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Message)
	}
	return c.Status(200).JSON(situs)
}

func (h *situsKeagamaanHandlerImpl) DeleteSitus(c *fiber.Ctx) error {
	situsId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.situsService.DeleteSitus(c.Context(), situsId); err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Message)
	}
	return c.SendStatus(204)
}

func (h *situsKeagamaanHandlerImpl) VerifySitus(c *fiber.Ctx) error {
	situsId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	var body dto.VerifikasiSitusRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	if err := h.situsService.VerifySitus(c.Context(), situsId, &body); err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Message)
	}
	return c.Status(200).SendString("Situs berhasil diverifikasi")
}

func (h *situsKeagamaanHandlerImpl) DeleteFotoSitus(c *fiber.Ctx) error {
	situsId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	var body dto.DeleteFoto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.validate.Struct(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	if err := h.situsService.DeleteFotoSitus(c.Context(), situsId, &body); err != nil {
		e := err.(*apperror.AppError)
		return c.Status(e.Status).SendString(e.Message)
	}
	return c.SendStatus(204)
}
