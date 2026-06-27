package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"modular-monolith/internal/modules/kelas/domain"
)

type pgKelasRepository struct {
	db *gorm.DB
}

func NewPgKelasRepository(db *gorm.DB) domain.KelasRepository {
	return &pgKelasRepository{db: db}
}

func (r *pgKelasRepository) Create(ctx context.Context, kelas *domain.Kelas) error {
	return r.db.WithContext(ctx).Create(kelas).Error
}

func (r *pgKelasRepository) GetAll(ctx context.Context) ([]*domain.Kelas, error) {
	var kelases []*domain.Kelas
	err := r.db.WithContext(ctx).Preload("ProgramStudi").Order("created_at desc").Find(&kelases).Error
	return kelases, err
}

func (r *pgKelasRepository) GetByID(ctx context.Context, id string) (*domain.Kelas, error) {
	var kelas domain.Kelas
	err := r.db.WithContext(ctx).Preload("ProgramStudi").Where("id = ?", id).First(&kelas).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kelas not found")
		}
		return nil, err
	}
	return &kelas, nil
}

func (r *pgKelasRepository) GetByName(ctx context.Context, name string) (*domain.Kelas, error) {
	var kelas domain.Kelas
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&kelas).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kelas not found")
		}
		return nil, err
	}
	return &kelas, nil
}

func (r *pgKelasRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Kelas{}).Error
}
