package domain

import (
	"context"
	"time"

	authDomain "modular-monolith/internal/modules/auth/domain"
	psDomain "modular-monolith/internal/modules/programstudi/domain"
)

type MataKuliah struct {
	ID             string                `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name           string                `json:"name" gorm:"type:varchar(255);uniqueIndex:idx_name_prodi;not null"`
	SKS            int                   `json:"sks" gorm:"not null"`
	ProgramStudiID string                 `json:"program_studi_id" gorm:"type:varchar(255);uniqueIndex:idx_name_prodi;not null"`
	ProgramStudi   *psDomain.ProgramStudi `json:"program_studi,omitempty" gorm:"foreignKey:ProgramStudiID"`
	Pengajuan      []*PengajuanMataKuliah `json:"pengajuan,omitempty" gorm:"foreignKey:MataKuliahID"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

type CreateMataKuliahRequest struct {
	Name           string `json:"name" validate:"required"`
	SKS            int    `json:"sks" validate:"required,min=1"`
	ProgramStudiID string `json:"program_studi_id" validate:"required"`
}

type MataKuliahRepository interface {
	Create(ctx context.Context, mk *MataKuliah) error
	GetByNameAndProdi(ctx context.Context, name string, prodiID string) (*MataKuliah, error)
	GetAll(ctx context.Context) ([]*MataKuliah, error)
	Delete(ctx context.Context, id string) error

	CreatePengajuan(ctx context.Context, p *PengajuanMataKuliah) error
	GetPengajuanByID(ctx context.Context, id string) (*PengajuanMataKuliah, error)
	GetPengajuanByDosenID(ctx context.Context, dosenID string) ([]*PengajuanMataKuliah, error)
	GetActivePengajuanByMataKuliahID(ctx context.Context, mkID string) ([]*PengajuanMataKuliah, error)
	GetAllPengajuan(ctx context.Context) ([]*PengajuanMataKuliah, error)
	UpdatePengajuan(ctx context.Context, p *PengajuanMataKuliah) error
}

type MataKuliahService interface {
	CreateMataKuliah(ctx context.Context, req CreateMataKuliahRequest) (*MataKuliah, error)
	GetMataKuliahList(ctx context.Context) ([]*MataKuliah, error)
	DeleteMataKuliah(ctx context.Context, id string) error

	RequestMataKuliah(ctx context.Context, dosenID string, req RequestMataKuliahPayload) (*PengajuanMataKuliah, error)
	ApprovePengajuan(ctx context.Context, id string) error
	RejectPengajuan(ctx context.Context, id string) error
	GetMyPengajuan(ctx context.Context, dosenID string) ([]*PengajuanMataKuliah, error)
	GetAllPengajuan(ctx context.Context) ([]*PengajuanMataKuliah, error)
}

type PengajuanMataKuliah struct {
	ID           string           `json:"id" gorm:"primaryKey;type:varchar(255)"`
	DosenID      string           `json:"dosen_id" gorm:"type:varchar(255);not null"`
	Dosen        *authDomain.User `json:"dosen,omitempty" gorm:"foreignKey:DosenID"`
	MataKuliahID string           `json:"mata_kuliah_id" gorm:"type:varchar(255);not null"`
	MataKuliah   *MataKuliah      `json:"mata_kuliah,omitempty" gorm:"foreignKey:MataKuliahID"`
	Status       string           `json:"status" gorm:"type:varchar(50);not null;default:'pending'"`
	Code         string           `json:"code" gorm:"type:varchar(6);not null"`
	CreatedAt    time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}

type RequestMataKuliahPayload struct {
	MataKuliahID string `json:"mata_kuliah_id" validate:"required"`
}

type ApproveRejectPengajuanPayload struct {
}
