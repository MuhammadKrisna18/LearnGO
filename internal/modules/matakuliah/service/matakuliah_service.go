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
	if req.ProgramStudiID == "" {
		return nil, apperrors.NewBadRequest("Program Studi wajib diisi")
	}

	// 1. Cek apakah mata kuliah dengan nama tersebut sudah ada (unik)
	existing, err := s.repo.GetByName(ctx, req.Name)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengecek mata kuliah", err.Error())
	}
	
	if existing != nil {
		return nil, &apperrors.AppError{Code: http.StatusConflict, Message: "Mata kuliah sudah ada, tidak bisa ditambahkan lagi meskipun berbeda SKS"}
	}

	newMk := &domain.MataKuliah{
		ID:             uuid.New().String(),
		Name:           req.Name,
		SKS:            req.SKS,
		ProgramStudiID: req.ProgramStudiID,
	}

	if err := s.repo.Create(ctx, newMk); err != nil {
		return nil, apperrors.NewInternal("Gagal menyimpan mata kuliah", err.Error())
	}

	return newMk, nil
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

func (s *matakuliahService) DeleteMataKuliah(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperrors.NewInternal("Gagal menghapus mata kuliah", err.Error())
	}
	return nil
}
