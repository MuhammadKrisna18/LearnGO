package service

import (
	"context"

	"modular-monolith/internal/modules/programstudi/domain"
	"modular-monolith/internal/shared/apperrors"
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
