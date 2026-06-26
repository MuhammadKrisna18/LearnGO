package domain

import (
	"context"
	"time"
)

type MataKuliah struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name      string    `json:"name" gorm:"type:varchar(255);uniqueIndex;not null"`
	SKS       int       `json:"sks" gorm:"type:int;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateMataKuliahRequest struct {
	Name string `json:"name" validate:"required"`
	SKS  int    `json:"sks" validate:"required,min=1"`
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
