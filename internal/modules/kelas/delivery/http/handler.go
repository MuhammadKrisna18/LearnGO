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

func (h *KelasHandler) RequestKelas(c *fiber.Ctx) error {
	dosenID := c.Locals("userID").(string)

	var req domain.RequestKelasPayload
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("Format payload tidak valid")
	}

	if err := h.validate.Struct(req); err != nil {
		return apperrors.NewBadRequest("Validasi gagal: kelas_id diperlukan")
	}

	pengajuan, err := h.service.RequestKelas(c.Context(), dosenID, req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusCreated, "Berhasil mengajukan kelas", pengajuan)
}

func (h *KelasHandler) GetMyPengajuan(c *fiber.Ctx) error {
	dosenID := c.Locals("userID").(string)

	list, err := h.service.GetMyPengajuan(c.Context(), dosenID)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil riwayat pengajuan kelas", list)
}

func (h *KelasHandler) GetAllPengajuan(c *fiber.Ctx) error {
	list, err := h.service.GetAllPengajuan(c.Context())
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil semua pengajuan kelas", list)
}

func (h *KelasHandler) ApprovePengajuan(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID pengajuan diperlukan")
	}

	if err := h.service.ApprovePengajuan(c.Context(), id); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil menyetujui kelas", nil)
}

func (h *KelasHandler) RejectPengajuan(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID pengajuan diperlukan")
	}

	if err := h.service.RejectPengajuan(c.Context(), id); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil menolak kelas", nil)
}
