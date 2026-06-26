package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
	"modular-monolith/config"
	_ "modular-monolith/docs"
	"modular-monolith/internal/modules/auth"
	"modular-monolith/internal/modules/matakuliah"
	"modular-monolith/internal/modules/programstudi"
	"modular-monolith/internal/shared/apperrors"
	"modular-monolith/internal/shared/cache"
	"modular-monolith/internal/shared/database"
	"modular-monolith/internal/shared/response"
	"modular-monolith/internal/middleware"
)

type App struct {
	cfg   *config.Config
	db    *gorm.DB
	redis *redis.Client
	fiber *fiber.App
}

func NewApp(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Start() error {
	var err error

	a.db, err = database.NewPostgresConnection(a.cfg)
	if err != nil {
		return fmt.Errorf("failed to init postgres: %w", err)
	}

	a.redis, err = cache.NewRedisClient(a.cfg)
	if err != nil {
		return fmt.Errorf("failed to init redis: %w", err)
	}

	a.fiber = fiber.New(fiber.Config{
		AppName:      "Modular Monolith API v1.0",
		ErrorHandler: globalErrorHandler,
	})

	a.fiber.Use(recover.New())
	a.fiber.Use(cors.New())
	
	a.fiber.Use(middleware.PerformanceLogger())
	
	a.fiber.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	a.fiber.Get("/", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, "System is healthy", fiber.Map{
			"env":   a.cfg.Env,
			"redis": "connected",
			"db":    "connected",
		})
	})

	a.fiber.Get("/swagger/*", swagger.HandlerDefault)

	api := a.fiber.Group("/api/v1")

	authModule := auth.NewAuthModule(a.db, a.cfg)
	authModule.RegisterRoutes(api)

	mkModule := matakuliah.NewMataKuliahModule(a.db, a.cfg)
	mkModule.RegisterRoutes(api)

	psModule := programstudi.NewProgramStudiModule(a.db, a.cfg)
	psModule.RegisterRoutes(api)

	serverErrors := make(chan error, 1)
	go func() {
		addr := fmt.Sprintf(":%d", a.cfg.Port)
		log.Printf("Server is running on %s\n", addr)
		serverErrors <- a.fiber.Listen(addr)
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Printf("Start shutdown, received signal: %v\n", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := a.fiber.ShutdownWithContext(ctx); err != nil {
			log.Printf("Could not gracefully shutdown server: %v\n", err)
		}

		if a.db != nil {
			log.Println("Closing PostgreSQL connection pool...")
			sqlDB, err := a.db.DB()
			if err == nil && sqlDB != nil {
				sqlDB.Close()
			}
		}

		if a.redis != nil {
			log.Println("Closing Redis client...")
			_ = a.redis.Close()
		}

		log.Println("Server gracefully stopped")
	}

	return nil
}

func globalErrorHandler(c *fiber.Ctx, err error) error {
	if e, ok := err.(*apperrors.AppError); ok {
		return response.Error(c, e.Code, e.Message, e.Details)
	}

	if e, ok := err.(*fiber.Error); ok {
		return response.Error(c, e.Code, e.Message, nil)
	}

	return response.Error(c, fiber.StatusInternalServerError, "An unexpected error occurred", err.Error())
}
