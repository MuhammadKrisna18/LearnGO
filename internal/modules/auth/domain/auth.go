package domain

import (
	"context"
	"time"

	psDomain "modular-monolith/internal/modules/programstudi/domain"
)

type User struct {
	ID             string                `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name           string                `json:"name" gorm:"type:varchar(255);not null"`
	Email          string                `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password       string                `json:"-" gorm:"type:varchar(255);not null"`
	Role           string                `json:"role" gorm:"type:varchar(50);not null;default:'user'"`
	ProgramStudiID *string               `json:"program_studi_id" gorm:"type:varchar(255)"`
	ProgramStudi   *psDomain.ProgramStudi `json:"program_studi,omitempty" gorm:"foreignKey:ProgramStudiID"`
	CreatedAt      time.Time             `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time             `json:"updated_at" gorm:"autoUpdateTime"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterDosenRequest struct {
	Name           string `json:"name" validate:"required"`
	Username       string `json:"username" validate:"required"`
	Password       string `json:"password" validate:"required"`
	ProgramStudiID string `json:"program_studi_id" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

type UserProfileResponse struct {
	ID             string                `json:"id"`
	Name           string                `json:"name"`
	Email          string                `json:"email"`
	Role           string                `json:"role"`
	ProgramStudiID *string               `json:"program_studi_id,omitempty"`
	ProgramStudi   *psDomain.ProgramStudi `json:"program_studi,omitempty"`
	CreatedAt      time.Time             `json:"created_at"`
}

type AuthRepository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetUsersByRole(ctx context.Context, role string) ([]*User, error)
	Create(ctx context.Context, user *User) error
}

type AuthService interface {
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	GetProfile(ctx context.Context, id string) (*UserProfileResponse, error)
	RegisterDosen(ctx context.Context, req RegisterDosenRequest) (*UserProfileResponse, error)
	GetDosenList(ctx context.Context) ([]*UserProfileResponse, error)
	DeleteDosen(ctx context.Context, id string) error
}
