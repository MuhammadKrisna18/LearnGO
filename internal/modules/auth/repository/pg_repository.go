package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"siakad-pro/internal/modules/auth/domain"
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
	result := r.db.WithContext(ctx).Preload("ProgramStudi").Where("id = ?", id).First(&u)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &u, nil
}

func (r *pgAuthRepository) GetUsersByRole(ctx context.Context, role string) ([]*domain.User, error) {
	var users []*domain.User
	result := r.db.WithContext(ctx).Preload("ProgramStudi").Where("role = ?", role).Order("created_at desc").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *pgAuthRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *pgAuthRepository) DeleteUser(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", id).Delete(&domain.EmailChangeRequest{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id).Delete(&domain.User{}).Error
	})
}

func (r *pgAuthRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *pgAuthRepository) CreateEmailChangeRequest(ctx context.Context, req *domain.EmailChangeRequest) error {
	return r.db.WithContext(ctx).Create(req).Error
}

func (r *pgAuthRepository) GetPendingEmailRequestByUserID(ctx context.Context, userID string) (*domain.EmailChangeRequest, error) {
	var req domain.EmailChangeRequest
	err := r.db.WithContext(ctx).Where("user_id = ? AND status = ?", userID, "pending").First(&req).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *pgAuthRepository) GetAllPendingEmailRequests(ctx context.Context) ([]*domain.EmailChangeRequest, error) {
	var reqs []*domain.EmailChangeRequest
	err := r.db.WithContext(ctx).Preload("User").Where("status = ?", "pending").Find(&reqs).Error
	return reqs, err
}

func (r *pgAuthRepository) GetEmailChangeRequestByID(ctx context.Context, id string) (*domain.EmailChangeRequest, error) {
	var req domain.EmailChangeRequest
	err := r.db.WithContext(ctx).Preload("User").Where("id = ?", id).First(&req).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *pgAuthRepository) UpdateEmailChangeRequest(ctx context.Context, req *domain.EmailChangeRequest) error {
	return r.db.WithContext(ctx).Save(req).Error
}
