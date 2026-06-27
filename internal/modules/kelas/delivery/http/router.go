package http

import (
	"github.com/gofiber/fiber/v2"
	"modular-monolith/internal/middleware"
)

func RegisterRoutes(router fiber.Router, handler *KelasHandler, jwtSecret string) {
	kelas := router.Group("/kelas")


	kelas.Use(middleware.Protected(jwtSecret))
	kelas.Get("/", handler.GetAll)


	kelas.Post("/", middleware.RequireRole("admin"), handler.Create)
	kelas.Delete("/:id", middleware.RequireRole("admin"), handler.Delete)
}
