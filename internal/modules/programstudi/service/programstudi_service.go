package service

import (
	"context"

	"siakad-pro/internal/modules/programstudi/domain"
	"siakad-pro/internal/shared/apperrors"
)

type programStudiService struct {
	repo domain.ProgramStudiRepository
}

func NewProgramStudiService(repo domain.ProgramStudiRepository) *programStudiService {
	return &programStudiService{repo: repo}
}

func (s *programStudiService) GetAll(ctx context.Context) ([]*domain.ProgramStudi, error) {
	psList, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengambil daftar program studi", err.Error())
	}
	return psList, nil
}
