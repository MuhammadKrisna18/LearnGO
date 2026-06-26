package service

import (
	"context"
	"net/http"
	"github.com/google/uuid"
	"modular-monolith/internal/modules/matakuliah/domain"
	"modular-monolith/internal/shared/apperrors"
)

type matakuliahService struct {
	repo domain.MataKuliahRepository
}

func NewMataKuliahService(repo domain.MataKuliahRepository) domain.MataKuliahService {
	return &matakuliahService{
		repo: repo,
	}
}

func (s *matakuliahService) CreateMataKuliah(ctx context.Context, req domain.CreateMataKuliahRequest) (*domain.MataKuliah, error) {
	// Check if already exists
	existing, err := s.repo.GetByName(ctx, req.Name)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengecek mata kuliah", err.Error())
	}
	
	if existing != nil {
		return nil, &apperrors.AppError{Code: http.StatusConflict, Message: "Mata kuliah sudah ada, tidak bisa ditambahkan lagi meskipun berbeda SKS"}
	}

	mk := &domain.MataKuliah{
		ID:   uuid.New().String(),
		Name: req.Name,
		SKS:  req.SKS,
	}

	if err := s.repo.Create(ctx, mk); err != nil {
		return nil, apperrors.NewInternal("Gagal menyimpan mata kuliah", err.Error())
	}

	return mk, nil
}

func (s *matakuliahService) GetMataKuliahList(ctx context.Context) ([]*domain.MataKuliah, error) {
	mkList, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengambil daftar mata kuliah", err.Error())
	}
	
	if mkList == nil {
		mkList = []*domain.MataKuliah{}
	}

	return mkList, nil
}
