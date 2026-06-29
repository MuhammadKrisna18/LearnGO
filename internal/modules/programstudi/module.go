package programstudi

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"modular-monolith/config"
	"modular-monolith/internal/modules/programstudi/delivery/http"
	"modular-monolith/internal/modules/programstudi/repository"
	"modular-monolith/internal/modules/programstudi/service"
)

type ProgramStudiModule struct {
	Handler *http.ProgramStudiHandler
}

func NewProgramStudiModule(db *gorm.DB, cfg *config.Config) *ProgramStudiModule {
	repo := repository.NewPgProgramStudiRepository(db)
	svc := service.NewProgramStudiService(repo)
	handler := http.NewProgramStudiHandler(svc)

	_ = repo.Seed(context.Background())

	return &ProgramStudiModule{
		Handler: handler,
	}
}

func (m *ProgramStudiModule) RegisterRoutes(router fiber.Router) {
	m.Handler.RegisterRoutes(router)
}
