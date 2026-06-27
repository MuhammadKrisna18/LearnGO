package service

import (
	"context"
	"regexp"

	"github.com/google/uuid"
	"modular-monolith/internal/modules/kelas/domain"
	"modular-monolith/internal/shared/apperrors"
)

type kelasService struct {
	repo domain.KelasRepository
}

func NewKelasService(repo domain.KelasRepository) domain.KelasService {
	return &kelasService{repo: repo}
}

func (s *kelasService) Create(ctx context.Context, req domain.CreateKelasRequest) (*domain.Kelas, error) {
	// Validate Name format: e.g. IF-101
	// The user requested format: IF-101 - IF 107, IF-201 - IF-207, IF-301
	// Regex: ^[A-Z0-9]+-\d{3}$
	matched, _ := regexp.MatchString(`^[A-Z0-9]+-\d{3}$`, req.Name)
	if !matched {
		return nil, apperrors.NewBadRequest("Format nama kelas tidak valid (contoh: IF-101)")
	}

	// Validate Capacity
	if req.Capacity < 25 || req.Capacity > 50 {
		return nil, apperrors.NewBadRequest("Kapasitas kelas harus antara 25 dan 50")
	}

	// Check for unique name
	existing, _ := s.repo.GetByName(ctx, req.Name)
	if existing != nil {
		return nil, apperrors.NewBadRequest("Nama kelas sudah terdaftar")
	}

	kelas := &domain.Kelas{
		ID:             uuid.New().String(),
		Name:           req.Name,
		Capacity:       req.Capacity,
		ProgramStudiID: req.ProgramStudiID,
	}

	if err := s.repo.Create(ctx, kelas); err != nil {
		return nil, apperrors.NewInternal("Gagal membuat kelas", err.Error())
	}

	return s.repo.GetByID(ctx, kelas.ID) // Fetch with relations
}

func (s *kelasService) GetAll(ctx context.Context) ([]*domain.Kelas, error) {
	return s.repo.GetAll(ctx)
}

func (s *kelasService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Kelas tidak ditemukan", err.Error())
	}
	
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperrors.NewInternal("Gagal menghapus kelas", err.Error())
	}
	return nil
}
