package domain

import (
	"context"
	"time"

	authDomain "modular-monolith/internal/modules/auth/domain"
	mkDomain "modular-monolith/internal/modules/matakuliah/domain"
	psDomain "modular-monolith/internal/modules/programstudi/domain"
)

type Kelas struct {
	ID             string                 `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name           string                 `json:"name" gorm:"type:varchar(255);not null"`
	Capacity       int                    `json:"capacity" gorm:"not null"`
	Hari           string                 `json:"hari" gorm:"type:varchar(20)"`
	JamMulai       string                 `json:"jam_mulai" gorm:"type:varchar(10)"`
	JamSelesai     string                 `json:"jam_selesai" gorm:"type:varchar(10)"`
	ProgramStudiID string                 `json:"program_studi_id" gorm:"type:varchar(255);not null"`
	ProgramStudi   *psDomain.ProgramStudi `json:"program_studi,omitempty" gorm:"foreignKey:ProgramStudiID"`
	Pengajuan      []*PengajuanKelas      `json:"pengajuan,omitempty" gorm:"foreignKey:KelasID"`
	CreatedAt      time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusRejected = "rejected"
)

type CreateKelasRequest struct {
	Name           string `json:"name" validate:"required"`
	Capacity       int    `json:"capacity" validate:"required,min=25,max=50"`
	Hari           string `json:"hari" validate:"required"`
	JamMulai       string `json:"jam_mulai" validate:"required"`
	JamSelesai     string `json:"jam_selesai" validate:"required"`
	ProgramStudiID string `json:"program_studi_id" validate:"required"`
}

type KelasRepository interface {
	Create(ctx context.Context, kelas *Kelas) error
	GetAll(ctx context.Context) ([]*Kelas, error)
	GetByID(ctx context.Context, id string) (*Kelas, error)
	GetByName(ctx context.Context, name string) (*Kelas, error)
	CheckScheduleConflict(ctx context.Context, name string, hari string, jamMulai string) (bool, error)
	Delete(ctx context.Context, id string) error
	HasApprovedMataKuliah(ctx context.Context, dosenID string) (bool, error)
	IsMataKuliahValidForKelas(ctx context.Context, dosenID string, mkID string, prodiID string) (bool, error)

	CreatePengajuan(ctx context.Context, p *PengajuanKelas) error
	GetPengajuanByID(ctx context.Context, id string) (*PengajuanKelas, error)
	GetPengajuanByDosenID(ctx context.Context, dosenID string) ([]*PengajuanKelas, error)
	GetActivePengajuanByKelasID(ctx context.Context, kelasID string) ([]*PengajuanKelas, error)
	GetAllPengajuan(ctx context.Context) ([]*PengajuanKelas, error)
	UpdatePengajuan(ctx context.Context, p *PengajuanKelas) error
	DeletePengajuan(ctx context.Context, id string) error
}

type KelasService interface {
	Create(ctx context.Context, req CreateKelasRequest) (*Kelas, error)
	GetAll(ctx context.Context) ([]*Kelas, error)
	Delete(ctx context.Context, id string) error

	RequestKelas(ctx context.Context, dosenID string, req RequestKelasPayload) (*PengajuanKelas, error)
	ApprovePengajuan(ctx context.Context, id string) error
	RejectPengajuan(ctx context.Context, id string) error
	GetMyPengajuan(ctx context.Context, dosenID string) ([]*PengajuanKelas, error)
	GetAllPengajuan(ctx context.Context) ([]*PengajuanKelas, error)
}

type PengajuanKelas struct {
	ID        string           `json:"id" gorm:"primaryKey;type:varchar(255)"`
	DosenID   string           `json:"dosen_id" gorm:"type:varchar(255);not null"`
	Dosen     *authDomain.User `json:"dosen,omitempty" gorm:"foreignKey:DosenID"`
	KelasID      string               `json:"kelas_id" gorm:"type:varchar(255);not null"`
	Kelas        *Kelas               `json:"kelas,omitempty" gorm:"foreignKey:KelasID"`
	MataKuliahID string               `json:"mata_kuliah_id" gorm:"type:varchar(255);not null"`
	MataKuliah   *mkDomain.MataKuliah `json:"mata_kuliah,omitempty" gorm:"foreignKey:MataKuliahID"`
	Status       string               `json:"status" gorm:"type:varchar(50);not null;default:'pending'"`
	Code         string               `json:"code" gorm:"type:varchar(6);not null"`
	CreatedAt time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}

type RequestKelasPayload struct {
	KelasID      string `json:"kelas_id" validate:"required"`
	MataKuliahID string `json:"mata_kuliah_id" validate:"required"`
}
