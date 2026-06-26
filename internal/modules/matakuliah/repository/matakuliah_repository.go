package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"modular-monolith/internal/modules/matakuliah/domain"
)

type pgMataKuliahRepository struct {
	db *gorm.DB
}

func NewPgMataKuliahRepository(db *gorm.DB) domain.MataKuliahRepository {
	return &pgMataKuliahRepository{db: db}
}

func (r *pgMataKuliahRepository) Create(ctx context.Context, mk *domain.MataKuliah) error {
	return r.db.WithContext(ctx).Create(mk).Error
}

func (r *pgMataKuliahRepository) GetByNameAndProdi(ctx context.Context, name string, prodiID string) (*domain.MataKuliah, error) {
	var mk domain.MataKuliah
	err := r.db.WithContext(ctx).Where("LOWER(name) = LOWER(?) AND program_studi_id = ?", name, prodiID).First(&mk).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &mk, nil
}

func (r *pgMataKuliahRepository) GetAll(ctx context.Context) ([]*domain.MataKuliah, error) {
	var mkList []*domain.MataKuliah
	err := r.db.WithContext(ctx).Preload("ProgramStudi").Order("created_at desc").Find(&mkList).Error
	if err != nil {
		return nil, err
	}
	return mkList, nil
}

func (r *pgMataKuliahRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.MataKuliah{}).Error
}
