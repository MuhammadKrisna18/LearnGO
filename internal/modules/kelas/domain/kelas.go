package domain

import (
	"context"
	"time"

	psDomain "modular-monolith/internal/modules/programstudi/domain"
)

type Kelas struct {
	ID             string                 `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name           string                 `json:"name" gorm:"type:varchar(255);uniqueIndex;not null"`
	Capacity       int                    `json:"capacity" gorm:"not null"`
	ProgramStudiID string                 `json:"program_studi_id" gorm:"type:varchar(255);not null"`
	ProgramStudi   *psDomain.ProgramStudi `json:"program_studi,omitempty" gorm:"foreignKey:ProgramStudiID"`
	CreatedAt      time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateKelasRequest struct {
	Name           string `json:"name" validate:"required,regexp=^[A-Z0-9]+-\\d{3}$"`
	Capacity       int    `json:"capacity" validate:"required,min=25,max=50"`
	ProgramStudiID string `json:"program_studi_id" validate:"required"`
}

type KelasRepository interface {
	Create(ctx context.Context, kelas *Kelas) error
	GetAll(ctx context.Context) ([]*Kelas, error)
	GetByID(ctx context.Context, id string) (*Kelas, error)
	GetByName(ctx context.Context, name string) (*Kelas, error)
	Delete(ctx context.Context, id string) error
}

type KelasService interface {
	Create(ctx context.Context, req CreateKelasRequest) (*Kelas, error)
	GetAll(ctx context.Context) ([]*Kelas, error)
	Delete(ctx context.Context, id string) error
}
