package kelas

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"siakad-pro/config"
	"siakad-pro/internal/modules/kelas/delivery/http"
	"siakad-pro/internal/modules/kelas/repository"
	"siakad-pro/internal/modules/kelas/service"
)

type KelasModule struct {
	handler *http.KelasHandler
	cfg     *config.Config
}

func NewKelasModule(db *gorm.DB, cfg *config.Config) *KelasModule {
	repo := repository.NewPgKelasRepository(db)
	svc := service.NewKelasService(repo)
	handler := http.NewKelasHandler(svc)

	return &KelasModule{
		handler: handler,
		cfg:     cfg,
	}
}

func (m *KelasModule) RegisterRoutes(router fiber.Router) {
	http.RegisterRoutes(router, m.handler, m.cfg.JWTSecret)
}
