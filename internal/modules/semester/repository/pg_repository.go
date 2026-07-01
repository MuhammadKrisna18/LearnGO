package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"siakad-pro/internal/modules/semester/domain"
)

type pgSemesterRepository struct {
	db *gorm.DB
}

func NewPgSemesterRepository(db *gorm.DB) domain.SemesterRepository {
	return &pgSemesterRepository{db: db}
}

func (r *pgSemesterRepository) Create(ctx context.Context, s *domain.Semester) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *pgSemesterRepository) GetAll(ctx context.Context) ([]*domain.Semester, error) {
	var semesters []*domain.Semester
	err := r.db.WithContext(ctx).Preload("MataKuliah").Preload("MataKuliah.MataKuliah").Preload("MataKuliah.MataKuliah.ProgramStudi").Order("nomor ASC").Find(&semesters).Error
	return semesters, err
}

func (r *pgSemesterRepository) GetByID(ctx context.Context, id string) (*domain.Semester, error) {
	var s domain.Semester
	err := r.db.WithContext(ctx).Preload("MataKuliah").Preload("MataKuliah.MataKuliah").Preload("MataKuliah.MataKuliah.ProgramStudi").First(&s, "id = ?", id).Error
	return &s, err
}

func (r *pgSemesterRepository) GetByNomor(ctx context.Context, nomor int) (*domain.Semester, error) {
	var s domain.Semester
	err := r.db.WithContext(ctx).First(&s, "nomor = ?", nomor).Error
	return &s, err
}

func (r *pgSemesterRepository) GetActive(ctx context.Context) (*domain.Semester, error) {
	var s domain.Semester
	err := r.db.WithContext(ctx).Preload("MataKuliah").Preload("MataKuliah.MataKuliah").First(&s, "is_active = ?", true).Error
	return &s, err
}

func (r *pgSemesterRepository) Update(ctx context.Context, s *domain.Semester) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *pgSemesterRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("semester_id = ?", id).Delete(&domain.SemesterMataKuliah{}).Error; err != nil {
			return err
		}
		if err := tx.Where("semester_id = ?", id).Delete(&domain.Pertemuan{}).Error; err != nil {
			return err
		}
		return tx.Delete(&domain.Semester{}, "id = ?", id).Error
	})
}

func (r *pgSemesterRepository) SetActive(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&domain.Semester{}).Where("id = ?", id).Update("is_active", true).Error
}

func (r *pgSemesterRepository) DeactivateAll(ctx context.Context) error {
	return r.db.WithContext(ctx).Model(&domain.Semester{}).Where("is_active = ?", true).Update("is_active", false).Error
}

func (r *pgSemesterRepository) AssignMataKuliah(ctx context.Context, sm *domain.SemesterMataKuliah) error {
	if sm.ID == "" {
		sm.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(sm).Error
}

func (r *pgSemesterRepository) UnassignMataKuliah(ctx context.Context, semesterID string, mkID string) error {
	return r.db.WithContext(ctx).Where("semester_id = ? AND mata_kuliah_id = ?", semesterID, mkID).Delete(&domain.SemesterMataKuliah{}).Error
}

func (r *pgSemesterRepository) GetSemesterMataKuliah(ctx context.Context, semesterID string) ([]*domain.SemesterMataKuliah, error) {
	var items []*domain.SemesterMataKuliah
	err := r.db.WithContext(ctx).Preload("MataKuliah").Preload("MataKuliah.ProgramStudi").Where("semester_id = ?", semesterID).Find(&items).Error
	return items, err
}

func (r *pgSemesterRepository) GetTotalSKS(ctx context.Context, semesterID string) (int, error) {
	var total int
	row := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(mk.sks), 0)
		FROM semester_mata_kuliahs smk
		JOIN mata_kuliahs mk ON mk.id = smk.mata_kuliah_id
		WHERE smk.semester_id = ?
	`, semesterID).Row()
	err := row.Scan(&total)
	return total, err
}

func (r *pgSemesterRepository) CreatePertemuan(ctx context.Context, p *domain.Pertemuan) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *pgSemesterRepository) GetPertemuanByKelasAndSemester(ctx context.Context, kelasID string, semesterID string) ([]*domain.Pertemuan, error) {
	var items []*domain.Pertemuan
	err := r.db.WithContext(ctx).Where("kelas_id = ? AND semester_id = ?", kelasID, semesterID).Order("nomor_pertemuan ASC").Find(&items).Error
	return items, err
}

func (r *pgSemesterRepository) CountSelesaiPertemuan(ctx context.Context, kelasID string, semesterID string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Pertemuan{}).Where("kelas_id = ? AND semester_id = ? AND status = ?", kelasID, semesterID, domain.PertemuanSelesai).Count(&count).Error
	return int(count), err
}

func (r *pgSemesterRepository) MarkPertemuanSelesai(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&domain.Pertemuan{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":          domain.PertemuanSelesai,
		"tanggal_selesai": now,
	}).Error
}

func (r *pgSemesterRepository) HasReachedMaxPertemuan(ctx context.Context, semesterID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Pertemuan{}).
		Where("semester_id = ? AND nomor_pertemuan >= ?", semesterID, domain.MaxPertemuan).
		Count(&count).Error
	return count > 0, err
}
