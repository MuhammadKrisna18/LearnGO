package semester

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"siakad-pro/config"
	"siakad-pro/internal/modules/semester/delivery/http"
	"siakad-pro/internal/modules/semester/repository"
	"siakad-pro/internal/modules/semester/service"
)

type SemesterModule struct {
	handler *http.SemesterHandler
	cfg     *config.Config
}

func NewSemesterModule(db *gorm.DB, cfg *config.Config) *SemesterModule {
	repo := repository.NewPgSemesterRepository(db)
	svc := service.NewSemesterService(repo)
	handler := http.NewSemesterHandler(svc)

	return &SemesterModule{
		handler: handler,
		cfg:     cfg,
	}
}

func (m *SemesterModule) RegisterRoutes(router fiber.Router) {
	http.RegisterRoutes(router, m.handler, m.cfg.JWTSecret)
}
