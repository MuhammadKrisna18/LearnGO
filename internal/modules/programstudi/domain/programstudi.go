package domain

import "context"

type ProgramStudi struct {
	ID   string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name string `json:"name" gorm:"type:varchar(255);uniqueIndex;not null"`
	Code string `json:"code" gorm:"type:varchar(50);uniqueIndex;not null"`
}

type ProgramStudiRepository interface {
	Create(ctx context.Context, ps *ProgramStudi) error
	GetByName(ctx context.Context, name string) (*ProgramStudi, error)
	GetAll(ctx context.Context) ([]*ProgramStudi, error)
	Seed(ctx context.Context) error
}
