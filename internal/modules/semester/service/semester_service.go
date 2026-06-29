package service

import (
	"context"
	"fmt"

	"siakad-pro/internal/modules/semester/domain"
	"siakad-pro/internal/shared/apperrors"
)

type semesterService struct {
	repo domain.SemesterRepository
}

func NewSemesterService(repo domain.SemesterRepository) domain.SemesterService {
	return &semesterService{repo: repo}
}

func (s *semesterService) Create(ctx context.Context, req domain.CreateSemesterRequest) (*domain.Semester, error) {
	if req.Nomor < domain.MinSemester || req.Nomor > domain.MaxSemester {
		return nil, apperrors.NewBadRequest(fmt.Sprintf("Nomor semester harus antara %d dan %d", domain.MinSemester, domain.MaxSemester))
	}

	if req.MinSKS >= req.MaxSKS {
		return nil, apperrors.NewBadRequest("Minimum SKS harus lebih kecil dari Maksimum SKS")
	}

	existing, _ := s.repo.GetByNomor(ctx, req.Nomor)
	if existing != nil {
		return nil, apperrors.NewConflict(fmt.Sprintf("Semester %d sudah ada", req.Nomor))
	}

	semester := &domain.Semester{
		Nomor:  req.Nomor,
		MinSKS: req.MinSKS,
		MaxSKS: req.MaxSKS,
	}

	if err := s.repo.Create(ctx, semester); err != nil {
		return nil, apperrors.NewInternal("Gagal membuat semester: " + err.Error())
	}

	return semester, nil
}

func (s *semesterService) GetAll(ctx context.Context) ([]*domain.Semester, error) {
	return s.repo.GetAll(ctx)
}

func (s *semesterService) GetByID(ctx context.Context, id string) (*domain.Semester, error) {
	sem, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFound("Semester tidak ditemukan")
	}
	return sem, nil
}

func (s *semesterService) Update(ctx context.Context, id string, req domain.UpdateSemesterRequest) (*domain.Semester, error) {
	sem, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFound("Semester tidak ditemukan")
	}

	if req.MinSKS != nil {
		sem.MinSKS = *req.MinSKS
	}
	if req.MaxSKS != nil {
		sem.MaxSKS = *req.MaxSKS
	}

	if sem.MinSKS >= sem.MaxSKS {
		return nil, apperrors.NewBadRequest("Minimum SKS harus lebih kecil dari Maksimum SKS")
	}

	if err := s.repo.Update(ctx, sem); err != nil {
		return nil, apperrors.NewInternal("Gagal mengupdate semester")
	}

	return sem, nil
}

func (s *semesterService) Delete(ctx context.Context, id string) error {
	sem, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Semester tidak ditemukan")
	}
	if sem.IsActive {
		return apperrors.NewBadRequest("Tidak bisa menghapus semester yang sedang aktif")
	}
	return s.repo.Delete(ctx, id)
}

func (s *semesterService) SetActive(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Semester tidak ditemukan")
	}

	if err := s.repo.DeactivateAll(ctx); err != nil {
		return apperrors.NewInternal("Gagal menonaktifkan semester sebelumnya")
	}

	if err := s.repo.SetActive(ctx, id); err != nil {
		return apperrors.NewInternal("Gagal mengaktifkan semester")
	}

	return nil
}

func (s *semesterService) AssignMataKuliah(ctx context.Context, semesterID string, req domain.AssignMataKuliahRequest) (*domain.SemesterMataKuliah, error) {
	sem, err := s.repo.GetByID(ctx, semesterID)
	if err != nil {
		return nil, apperrors.NewNotFound("Semester tidak ditemukan")
	}

	for _, smk := range sem.MataKuliah {
		if smk.MataKuliahID == req.MataKuliahID {
			return nil, apperrors.NewConflict("Mata kuliah sudah ada di semester ini")
		}
	}

	totalSKS, err := s.repo.GetTotalSKS(ctx, semesterID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal menghitung total SKS")
	}

	_ = totalSKS

	sm := &domain.SemesterMataKuliah{
		SemesterID:   semesterID,
		MataKuliahID: req.MataKuliahID,
	}

	if err := s.repo.AssignMataKuliah(ctx, sm); err != nil {
		return nil, apperrors.NewInternal("Gagal menambahkan mata kuliah ke semester")
	}

	return sm, nil
}

func (s *semesterService) UnassignMataKuliah(ctx context.Context, semesterID string, mkID string) error {
	return s.repo.UnassignMataKuliah(ctx, semesterID, mkID)
}

func (s *semesterService) CatatPertemuan(ctx context.Context, dosenID string, req domain.CatatPertemuanRequest) (*domain.Pertemuan, error) {
	existing, err := s.repo.GetPertemuanByKelasAndSemester(ctx, req.KelasID, req.SemesterID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengambil data pertemuan")
	}

	nomorNext := len(existing) + 1
	if nomorNext > domain.MaxPertemuan {
		return nil, apperrors.NewBadRequest("Semua 16 pertemuan sudah tercatat untuk kelas ini")
	}

	p := &domain.Pertemuan{
		SemesterID:     req.SemesterID,
		KelasID:        req.KelasID,
		DosenID:        dosenID,
		NomorPertemuan: nomorNext,
		Topik:          req.Topik,
		Status:         domain.PertemuanSelesai,
	}

	if err := s.repo.CreatePertemuan(ctx, p); err != nil {
		return nil, apperrors.NewInternal("Gagal mencatat pertemuan")
	}

	if nomorNext == domain.MaxPertemuan {
		sem, _ := s.repo.GetByID(ctx, req.SemesterID)
		if sem != nil && sem.IsActive && sem.Nomor < domain.MaxSemester {
			nextSem, _ := s.repo.GetByNomor(ctx, sem.Nomor+1)
			if nextSem != nil {
				_ = s.repo.DeactivateAll(ctx)
				_ = s.repo.SetActive(ctx, nextSem.ID)
			}
		}
	}

	return p, nil
}

func (s *semesterService) GetPertemuan(ctx context.Context, kelasID string, semesterID string) ([]*domain.Pertemuan, error) {
	return s.repo.GetPertemuanByKelasAndSemester(ctx, kelasID, semesterID)
}
