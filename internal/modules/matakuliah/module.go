package matakuliah

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"modular-monolith/config"
	"modular-monolith/internal/modules/matakuliah/delivery/http"
	"modular-monolith/internal/modules/matakuliah/repository"
	"modular-monolith/internal/modules/matakuliah/service"
)

type Module struct {
	Handler *http.MataKuliahHandler
	cfg     *config.Config
}

func NewMataKuliahModule(db *gorm.DB, cfg *config.Config) *Module {
	repo := repository.NewPgMataKuliahRepository(db)
	svc := service.NewMataKuliahService(repo)
	handler := http.NewMataKuliahHandler(svc)

	return &Module{
		Handler: handler,
		cfg:     cfg,
	}
}

func (m *Module) RegisterRoutes(router fiber.Router) {
	m.Handler.RegisterRoutes(router, m.cfg.JWTSecret)
}
