package http

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"modular-monolith/internal/modules/programstudi/domain"
	"modular-monolith/internal/shared/response"
)

type ProgramStudiHandler struct {
	service *ServiceWrapper
}

type ProgramStudiServiceInterface interface {
	GetAll(ctx context.Context) ([]*domain.ProgramStudi, error)
}

func NewProgramStudiHandler(service ProgramStudiServiceInterface) *ProgramStudiHandler {
	return &ProgramStudiHandler{service: &ServiceWrapper{svc: service}}
}

type ServiceWrapper struct {
	svc ProgramStudiServiceInterface
}

func (h *ProgramStudiHandler) RegisterRoutes(router fiber.Router) {
	psGroup := router.Group("/program-studi")
	psGroup.Get("/", h.GetAll)
}

func (h *ProgramStudiHandler) GetAll(c *fiber.Ctx) error {
	list, err := h.service.svc.GetAll(c.UserContext())
	if err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Berhasil mengambil daftar program studi", list)
}
