package http

import (
	"github.com/gofiber/fiber/v2"
	"siakad-pro/internal/modules/semester/domain"
	"siakad-pro/internal/shared/apperrors"
	"siakad-pro/internal/shared/response"
)

type SemesterHandler struct {
	service domain.SemesterService
}

func NewSemesterHandler(svc domain.SemesterService) *SemesterHandler {
	return &SemesterHandler{service: svc}
}

func (h *SemesterHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateSemesterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	sem, err := h.service.Create(c.Context(), req)
	if err != nil {
		if e, ok := err.(*apperrors.AppError); ok {
			return response.Error(c, e.Code, e.Message, nil)
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusCreated, "Semester berhasil dibuat", sem)
}

func (h *SemesterHandler) GetAll(c *fiber.Ctx) error {
	semesters, err := h.service.GetAll(c.Context())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, fiber.StatusOK, "Daftar semester", semesters)
}

func (h *SemesterHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	sem, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		if e, ok := err.(*apperrors.AppError); ok {
			return response.Error(c, e.Code, e.Message, nil)
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, fiber.StatusOK, "Detail semester", sem)
}

func (h *SemesterHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req domain.UpdateSemesterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	sem, err := h.service.Update(c.Context(), id, req)
	if err != nil {
		if e, ok := err.(*apperrors.AppError); ok {
			return response.Error(c, e.Code, e.Message, nil)
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Semester berhasil diupdate", sem)
}

func (h *SemesterHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.Delete(c.Context(), id); err != nil {
		if e, ok := err.(*apperrors.AppError); ok {
			return response.Error(c, e.Code, e.Message, nil)
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, fiber.StatusOK, "Semester berhasil dihapus", nil)
}

func (h *SemesterHandler) SetActive(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.SetActive(c.Context(), id); err != nil {
		if e, ok := err.(*apperrors.AppError); ok {
			return response.Error(c, e.Code, e.Message, nil)
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, fiber.StatusOK, "Semester berhasil diaktifkan", nil)
}

func (h *SemesterHandler) AssignMataKuliah(c *fiber.Ctx) error {
	semesterID := c.Params("id")
	var req domain.AssignMataKuliahRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	sm, err := h.service.AssignMataKuliah(c.Context(), semesterID, req)
	if err != nil {
		if e, ok := err.(*apperrors.AppError); ok {
			return response.Error(c, e.Code, e.Message, nil)
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusCreated, "Mata kuliah berhasil ditambahkan ke semester", sm)
}

func (h *SemesterHandler) UnassignMataKuliah(c *fiber.Ctx) error {
	semesterID := c.Params("id")
	mkID := c.Params("mkId")
	if err := h.service.UnassignMataKuliah(c.Context(), semesterID, mkID); err != nil {
		if e, ok := err.(*apperrors.AppError); ok {
			return response.Error(c, e.Code, e.Message, nil)
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, fiber.StatusOK, "Mata kuliah berhasil dihapus dari semester", nil)
}

func (h *SemesterHandler) CatatPertemuan(c *fiber.Ctx) error {
	dosenID := c.Locals("userID").(string)
	var req domain.CatatPertemuanRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	p, err := h.service.CatatPertemuan(c.Context(), dosenID, req)
	if err != nil {
		if e, ok := err.(*apperrors.AppError); ok {
			return response.Error(c, e.Code, e.Message, nil)
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusCreated, "Pertemuan berhasil dicatat", p)
}

func (h *SemesterHandler) GetPertemuan(c *fiber.Ctx) error {
	kelasID := c.Query("kelas_id")
	semesterID := c.Query("semester_id")

	if kelasID == "" || semesterID == "" {
		return response.Error(c, fiber.StatusBadRequest, "Parameter kelas_id dan semester_id diperlukan", nil)
	}

	items, err := h.service.GetPertemuan(c.Context(), kelasID, semesterID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Daftar pertemuan", items)
}
