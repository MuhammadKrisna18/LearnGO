package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"modular-monolith/internal/modules/kelas/domain"
	"modular-monolith/internal/shared/apperrors"
	"modular-monolith/internal/shared/response"
)

type KelasHandler struct {
	service  domain.KelasService
	validate *validator.Validate
}

func NewKelasHandler(service domain.KelasService) *KelasHandler {
	return &KelasHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *KelasHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateKelasRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("Format payload tidak valid")
	}

	if err := h.validate.Struct(req); err != nil {
		return apperrors.NewBadRequest("Validasi gagal: periksa kembali data yang dimasukkan")
	}

	kelas, err := h.service.Create(c.Context(), req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusCreated, "Kelas berhasil ditambahkan", kelas)
}

func (h *KelasHandler) GetAll(c *fiber.Ctx) error {
	kelases, err := h.service.GetAll(c.Context())
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil data kelas", kelases)
}

func (h *KelasHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID kelas diperlukan")
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Kelas berhasil dihapus", nil)
}
