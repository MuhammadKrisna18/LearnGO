package service

import (
	"context"
	"log"
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

	matched, _ := regexp.MatchString(`^IF-[1-3]0[1-7]$`, req.Name)
	if !matched {
		return nil, apperrors.NewBadRequest("Format nama kelas tidak valid (contoh yang benar: IF-101 s/d IF-107, IF-201 s/d IF-207, IF-301 s/d IF-307)")
	}


	if req.Capacity < 25 || req.Capacity > 50 {
		return nil, apperrors.NewBadRequest("Kapasitas kelas harus antara 25 dan 50")
	}


	conflict, _ := s.repo.CheckScheduleConflict(ctx, req.Name, req.Hari, req.JamMulai)
	if conflict {
		return nil, apperrors.NewBadRequest("Kelas tersebut sudah terdaftar pada hari dan jam yang sama")
	}


	validHari := map[string]bool{"Senin": true, "Selasa": true, "Rabu": true, "Kamis": true, "Jumat": true}
	if !validHari[req.Hari] {
		return nil, apperrors.NewBadRequest("Hari harus antara Senin sampai Jumat")
	}


	if req.JamMulai == "" || req.JamSelesai == "" {
		return nil, apperrors.NewBadRequest("Jam mulai dan selesai harus diisi")
	}

	kelas := &domain.Kelas{
		ID:             uuid.New().String(),
		Name:           req.Name,
		Capacity:       req.Capacity,
		Hari:           req.Hari,
		JamMulai:       req.JamMulai,
		JamSelesai:     req.JamSelesai,
		ProgramStudiID: req.ProgramStudiID,
	}

	if err := s.repo.Create(ctx, kelas); err != nil {
		log.Println("DB INSERT ERROR:", err)
		return nil, apperrors.NewInternal("Gagal membuat kelas: " + err.Error())
	}

	return s.repo.GetByID(ctx, kelas.ID)
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
