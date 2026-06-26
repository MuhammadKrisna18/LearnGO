package domain

import (
	"context"
	"time"

	psDomain "modular-monolith/internal/modules/programstudi/domain"
)

type MataKuliah struct {
	ID             string                `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name           string                `json:"name" gorm:"type:varchar(255);uniqueIndex;not null"`
	SKS            int                   `json:"sks" gorm:"not null"`
	ProgramStudiID string                `json:"program_studi_id" gorm:"type:varchar(255);not null"`
	ProgramStudi   *psDomain.ProgramStudi `json:"program_studi,omitempty" gorm:"foreignKey:ProgramStudiID"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
}

type CreateMataKuliahRequest struct {
	Name           string `json:"name" validate:"required"`
	SKS            int    `json:"sks" validate:"required,min=1"`
	ProgramStudiID string `json:"program_studi_id" validate:"required"`
}

type MataKuliahRepository interface {
	Create(ctx context.Context, mk *MataKuliah) error
	GetByName(ctx context.Context, name string) (*MataKuliah, error)
	GetAll(ctx context.Context) ([]*MataKuliah, error)
	Delete(ctx context.Context, id string) error
}

type MataKuliahService interface {
	CreateMataKuliah(ctx context.Context, req CreateMataKuliahRequest) (*MataKuliah, error)
	GetMataKuliahList(ctx context.Context) ([]*MataKuliah, error)
	DeleteMataKuliah(ctx context.Context, id string) error
}
