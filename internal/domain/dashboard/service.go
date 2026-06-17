package dashboard

import "go-library/internal/domain/dashboard/dto"

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) GetStats() (*dto.Response, error) {
	return s.repo.GetStats()
}

func (s *service) GetUserStats(userID uint) (*dto.UserStats, error) {
	return s.repo.GetUserStats(userID)
}
