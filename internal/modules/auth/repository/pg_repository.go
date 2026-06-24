package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"modular-monolith/internal/modules/auth/domain"
)

type pgAuthRepository struct {
	db *gorm.DB
}

func NewPgAuthRepository(db *gorm.DB) domain.AuthRepository {
	return &pgAuthRepository{db: db}
}

func (r *pgAuthRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&u)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &u, nil
}

func (r *pgAuthRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&u)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &u, nil
}

func (r *pgAuthRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
