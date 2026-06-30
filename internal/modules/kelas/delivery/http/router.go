package http

import (
	"github.com/gofiber/fiber/v2"
	"siakad-pro/internal/middleware"
)

func RegisterRoutes(router fiber.Router, handler *KelasHandler, jwtSecret string) {
	kelas := router.Group("/kelas")

	kelas.Use(middleware.Protected(jwtSecret))
	kelas.Get("/", handler.GetAll)

	kelas.Post("/request", middleware.RequireRole("dosen"), handler.RequestKelas)
	kelas.Get("/pengajuan/me", middleware.RequireRole("dosen"), handler.GetMyPengajuan)
	kelas.Get("/pengajuan/:id/mahasiswa", middleware.RequireRole("dosen"), handler.GetMahasiswaInKelas)

	kelas.Get("/mahasiswa/my-jadwal", middleware.RequireRole("mahasiswa"), handler.GetMyJadwal)

	kelas.Get("/pengajuan", middleware.RequireRole("admin"), handler.GetAllPengajuan)
	kelas.Post("/pengajuan/:id/approve", middleware.RequireRole("admin"), handler.ApprovePengajuan)
	kelas.Post("/pengajuan/:id/reject", middleware.RequireRole("admin"), handler.RejectPengajuan)

	kelas.Post("/", middleware.RequireRole("admin"), handler.Create)
	kelas.Delete("/:id", middleware.RequireRole("admin"), handler.Delete)
	kelas.Get("/:id", handler.GetByID)
}
