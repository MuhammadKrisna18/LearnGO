package database

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"modular-monolith/config"
	"modular-monolith/internal/modules/auth/domain"
	mkDomain "modular-monolith/internal/modules/matakuliah/domain"
	psDomain "modular-monolith/internal/modules/programstudi/domain"
)

func NewPostgresConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Error),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("unable to get underlying db: %w", err)
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	log.Println("Successfully connected to PostgreSQL database via GORM")

	log.Println("Running AutoMigrate...")
	if err := db.AutoMigrate(
		&domain.User{},
		&mkDomain.MataKuliah{},
		&psDomain.ProgramStudi{},
		&domain.EmailChangeRequest{},
	); err != nil {
		log.Printf("AutoMigrate failed: %v", err)
	}

	seedAdmin(db)

	return db, nil
}

func seedAdmin(db *gorm.DB) {
	adminEmail := "adminGO@golang.id"
	var count int64
	db.Model(&domain.User{}).Where("email = ?", adminEmail).Count(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("adminGO"), bcrypt.DefaultCost)
		adminUser := domain.User{
			ID:       "admin-id-seeder-12345",
			Name:     "Super Admin",
			Email:    adminEmail,
			Password: string(hashedPassword),
			Role:     "admin",
		}
		if err := db.Create(&adminUser).Error; err != nil {
			log.Printf("Failed to seed admin user: %v", err)
		} else {
			log.Println("Successfully seeded admin user: adminGO@golang.id")
		}
	}
}
