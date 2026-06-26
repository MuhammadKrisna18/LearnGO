package http

import (
	"github.com/gofiber/fiber/v2"
	"modular-monolith/internal/middleware"
	"modular-monolith/internal/modules/matakuliah/domain"
	"modular-monolith/internal/shared/response"
)

type MataKuliahHandler struct {
	service domain.MataKuliahService
}

func NewMataKuliahHandler(service domain.MataKuliahService) *MataKuliahHandler {
	return &MataKuliahHandler{
		service: service,
	}
}

func (h *MataKuliahHandler) RegisterRoutes(router fiber.Router, jwtSecret string) {
	mkGroup := router.Group("/matakuliah")

	mkGroup.Use(middleware.Protected(jwtSecret))
	
	mkGroup.Get("/", h.GetMataKuliahList)
	mkGroup.Post("/", middleware.RequireRole("admin"), h.CreateMataKuliah)
}

func (h *MataKuliahHandler) CreateMataKuliah(c *fiber.Ctx) error {
	var req domain.CreateMataKuliahRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload", err.Error())
	}

	if req.Name == "" || req.SKS < 1 {
		return response.Error(c, fiber.StatusBadRequest, "Nama mata kuliah dan SKS harus diisi", nil)
	}

	mk, err := h.service.CreateMataKuliah(c.Context(), req)
	if err != nil {
		return err // handled by global error handler
	}

	return response.Success(c, fiber.StatusCreated, "Mata kuliah berhasil ditambahkan", mk)
}

func (h *MataKuliahHandler) GetMataKuliahList(c *fiber.Ctx) error {
	mkList, err := h.service.GetMataKuliahList(c.Context())
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil daftar mata kuliah", mkList)
}
