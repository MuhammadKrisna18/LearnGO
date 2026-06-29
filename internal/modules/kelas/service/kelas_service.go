package service

import (
	"context"
	"log"
	"regexp"
	"time"

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
		return apperrors.NewNotFound("Kelas tidak ditemukan")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		log.Printf("Error deleting kelas: %v", err)
		return apperrors.NewInternal("Gagal menghapus kelas", err.Error())
	}
	return nil
}

// RequestKelas
func (s *kelasService) RequestKelas(ctx context.Context, dosenID string, req domain.RequestKelasPayload) (*domain.PengajuanKelas, error) {
	kelas, err := s.repo.GetByID(ctx, req.KelasID)
	if err != nil {
		return nil, apperrors.NewNotFound("Kelas tidak ditemukan")
	}

	hasMK, err := s.repo.IsMataKuliahValidForKelas(ctx, dosenID, req.MataKuliahID, kelas.ProgramStudiID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal memvalidasi mata kuliah dosen", err.Error())
	}
	if !hasMK {
		return nil, apperrors.NewBadRequest("Mata kuliah yang dipilih tidak valid atau bukan berasal dari program studi kelas ini.")
	}

	activePengajuan, err := s.repo.GetActivePengajuanByKelasID(ctx, req.KelasID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengecek status kelas", err.Error())
	}

	for _, p := range activePengajuan {
		if p.Status == domain.StatusApproved {
			return nil, apperrors.NewBadRequest("Kelas ini sudah disetujui untuk dosen lain")
		}
		if p.DosenID == dosenID && p.Status == domain.StatusPending {
			return nil, apperrors.NewBadRequest("Anda sudah mengajukan kelas ini")
		}
	}

	// Validasi bentrok jadwal dosen
	dosenPengajuan, err := s.repo.GetPengajuanByDosenID(ctx, dosenID)
	if err == nil {
		for _, dp := range dosenPengajuan {
			if dp.Status != domain.StatusRejected && dp.Kelas != nil {
				// Cek bentrok hari dan jam
				if dp.Kelas.Hari == kelas.Hari {
					// Pengecekan overlap waktu (format HH:MM)
					// Overlap terjadi jika: Mulai1 < Selesai2 && Mulai2 < Selesai1
					if kelas.JamMulai < dp.Kelas.JamSelesai && dp.Kelas.JamMulai < kelas.JamSelesai {
						return nil, apperrors.NewBadRequest("Jadwal kelas bentrok dengan kelas lain yang sudah Anda ajukan/ambil (" + dp.Kelas.Name + ")")
					}
				}
			}
		}
	}

	pengajuan := &domain.PengajuanKelas{
		ID:           uuid.NewString(),
		DosenID:      dosenID,
		KelasID:      req.KelasID,
		MataKuliahID: req.MataKuliahID,
		Status:       domain.StatusPending,
		Code:         generateRandomCode(),
	}

	if err := s.repo.CreatePengajuan(ctx, pengajuan); err != nil {
		return nil, apperrors.NewInternal("Gagal mengajukan kelas", err.Error())
	}

	pengajuan.Kelas = kelas
	return pengajuan, nil
}

func (s *kelasService) ApprovePengajuan(ctx context.Context, id string) error {
	p, err := s.repo.GetPengajuanByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Pengajuan tidak ditemukan")
	}

	if p.Status != domain.StatusPending {
		return apperrors.NewBadRequest("Hanya pengajuan berstatus pending yang dapat disetujui")
	}

	// Cek apakah sudah ada yang di-approve untuk kelas yang sama
	activePengajuan, err := s.repo.GetActivePengajuanByKelasID(ctx, p.KelasID)
	if err != nil {
		return apperrors.NewInternal("Gagal mengecek status kelas", err.Error())
	}

	for _, ap := range activePengajuan {
		if ap.Status == domain.StatusApproved && ap.ID != id {
			return apperrors.NewBadRequest("Kelas sudah disetujui untuk dosen lain")
		}
	}

	p.Status = domain.StatusApproved
	if err := s.repo.UpdatePengajuan(ctx, p); err != nil {
		return apperrors.NewInternal("Gagal menyetujui pengajuan", err.Error())
	}

	// Reject others
	for _, ap := range activePengajuan {
		if ap.ID != id && ap.Status == domain.StatusPending {
			ap.Status = domain.StatusRejected
			s.repo.UpdatePengajuan(ctx, ap)
		}
	}

	return nil
}

func (s *kelasService) RejectPengajuan(ctx context.Context, id string) error {
	p, err := s.repo.GetPengajuanByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Pengajuan tidak ditemukan")
	}

	if p.Status != domain.StatusPending {
		return apperrors.NewBadRequest("Hanya pengajuan berstatus pending yang dapat ditolak")
	}

	p.Status = domain.StatusRejected
	if err := s.repo.UpdatePengajuan(ctx, p); err != nil {
		return apperrors.NewInternal("Gagal menolak pengajuan", err.Error())
	}

	return nil
}

func (s *kelasService) GetMyPengajuan(ctx context.Context, dosenID string) ([]*domain.PengajuanKelas, error) {
	return s.repo.GetPengajuanByDosenID(ctx, dosenID)
}

func (s *kelasService) GetAllPengajuan(ctx context.Context) ([]*domain.PengajuanKelas, error) {
	return s.repo.GetAllPengajuan(ctx)
}

// Helper to generate a random 6-character alphanumeric code
func generateRandomCode() string {
	b := make([]byte, 6)
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	importRand := true
	for i := range b {
		b[i] = charset[fastRand()%len(charset)]
	}
	_ = importRand
	return string(b)
}

var fastRand = func() int {
	// using time.Now().UnixNano() directly to avoid importing math/rand or crypto/rand for simplicity if not already imported.
	// We'll use a simple approach with crypto/rand for real scenarios, but we can import math/rand locally.
	return int(time.Now().UnixNano())
}
