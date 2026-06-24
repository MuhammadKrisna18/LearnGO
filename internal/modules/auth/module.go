package auth

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"modular-monolith/config"
	"modular-monolith/internal/modules/auth/delivery/http"
	"modular-monolith/internal/modules/auth/repository"
	"modular-monolith/internal/modules/auth/service"
)

type Module struct {
	Handler *http.AuthHandler
	cfg     *config.Config
}

func NewAuthModule(db *gorm.DB, cfg *config.Config) *Module {
	repo := repository.NewPgAuthRepository(db)
	svc := service.NewAuthService(repo, cfg)
	handler := http.NewAuthHandler(svc)

	return &Module{
		Handler: handler,
		cfg:     cfg,
	}
}

func (m *Module) RegisterRoutes(router fiber.Router) {
	m.Handler.RegisterRoutes(router, m.cfg.JWTSecret)
}
