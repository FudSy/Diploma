package service

import (
	"github.com/FudSy/Diploma/internal/dto"
	"github.com/FudSy/Diploma/internal/pkg/repository"
)

type StatsService struct {
	statsRepo repository.Stats
}

func NewStatsService(r repository.Stats) *StatsService {
	return &StatsService{statsRepo: r}
}

func (s *StatsService) GetOverview() (dto.StatsOverview, error) {
	return s.statsRepo.GetOverview()
}
