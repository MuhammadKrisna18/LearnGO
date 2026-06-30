package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	authDomain "siakad-pro/internal/modules/auth/domain"
	"siakad-pro/internal/modules/kelas/domain"
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
	err := r.db.WithContext(ctx).
		Preload("ProgramStudi").
		Preload("Pengajuan", "status IN ?", []string{domain.StatusPending, domain.StatusApproved}).
		Preload("Pengajuan.Dosen").
		Order("created_at desc").Find(&kelases).Error
	return kelases, err
}

func (r *pgKelasRepository) GetByID(ctx context.Context, id string) (*domain.Kelas, error) {
	var kelas domain.Kelas
	err := r.db.WithContext(ctx).
		Preload("ProgramStudi").
		Preload("Pengajuan", "status IN ?", []string{domain.StatusPending, domain.StatusApproved}).
		Preload("Pengajuan.Dosen").
		Preload("Pengajuan.MataKuliah").
		Where("id = ?", id).First(&kelas).Error
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

func (r *pgKelasRepository) CheckScheduleConflict(ctx context.Context, name string, hari string, jamMulai string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Kelas{}).
		Where("name = ? AND hari = ? AND jam_mulai = ?", name, hari, jamMulai).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *pgKelasRepository) HasApprovedMataKuliah(ctx context.Context, dosenID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("pengajuan_mata_kuliahs").
		Where("dosen_id = ? AND status = ?", dosenID, domain.StatusApproved).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *pgKelasRepository) IsMataKuliahValidForKelas(ctx context.Context, dosenID string, mkID string, prodiID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("pengajuan_mata_kuliahs").
		Joins("JOIN mata_kuliahs mk ON mk.id = pengajuan_mata_kuliahs.mata_kuliah_id").
		Where("pengajuan_mata_kuliahs.dosen_id = ? AND pengajuan_mata_kuliahs.mata_kuliah_id = ? AND pengajuan_mata_kuliahs.status = ? AND mk.program_studi_id = ?", dosenID, mkID, domain.StatusApproved, prodiID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *pgKelasRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Kelas{}).Error
}

func (r *pgKelasRepository) CreatePengajuan(ctx context.Context, p *domain.PengajuanKelas) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *pgKelasRepository) GetPengajuanByID(ctx context.Context, id string) (*domain.PengajuanKelas, error) {
	var p domain.PengajuanKelas
	err := r.db.WithContext(ctx).Preload("Kelas").Preload("Dosen").Preload("MataKuliah").Where("id = ?", id).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengajuan not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *pgKelasRepository) GetPengajuanByDosenID(ctx context.Context, dosenID string) ([]*domain.PengajuanKelas, error) {
	var list []*domain.PengajuanKelas
	err := r.db.WithContext(ctx).Preload("Kelas").Preload("Kelas.ProgramStudi").Preload("MataKuliah").Where("dosen_id = ?", dosenID).Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *pgKelasRepository) GetActivePengajuanByKelasID(ctx context.Context, kelasID string) ([]*domain.PengajuanKelas, error) {
	var list []*domain.PengajuanKelas
	err := r.db.WithContext(ctx).Where("kelas_id = ? AND status IN ?", kelasID, []string{domain.StatusPending, domain.StatusApproved}).Find(&list).Error
	return list, err
}

func (r *pgKelasRepository) GetAllPengajuan(ctx context.Context) ([]*domain.PengajuanKelas, error) {
	var list []*domain.PengajuanKelas
	err := r.db.WithContext(ctx).Preload("Kelas").Preload("Kelas.ProgramStudi").Preload("Dosen").Preload("MataKuliah").Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *pgKelasRepository) UpdatePengajuan(ctx context.Context, p *domain.PengajuanKelas) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *pgKelasRepository) DeletePengajuan(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.PengajuanKelas{}).Error
}

func (r *pgKelasRepository) GetMahasiswaByProgramStudiID(ctx context.Context, prodiID string) ([]*authDomain.User, error) {
	var users []*authDomain.User
	err := r.db.WithContext(ctx).Table("users").Where("role = ? AND program_studi_id = ?", "mahasiswa", prodiID).Find(&users).Error
	return users, err
}

func (r *pgKelasRepository) GetApprovedPengajuanByProdiID(ctx context.Context, prodiID string) ([]*domain.PengajuanKelas, error) {
	var list []*domain.PengajuanKelas
	err := r.db.WithContext(ctx).
		Joins("JOIN kelas ON kelas.id = pengajuan_kelas.kelas_id").
		Preload("Kelas").
		Preload("Dosen").
		Preload("MataKuliah").
		Where("kelas.program_studi_id = ? AND pengajuan_kelas.status = ?", prodiID, domain.StatusApproved).
		Order("pengajuan_kelas.created_at desc").
		Find(&list).Error
	return list, err
}

func (r *pgKelasRepository) GetUserByID(ctx context.Context, userID string) (*authDomain.User, error) {
	var user authDomain.User
	err := r.db.WithContext(ctx).Table("users").Where("id = ?", userID).First(&user).Error
	return &user, err
}
