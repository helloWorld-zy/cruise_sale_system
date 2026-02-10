package core

import (
	"context"
)

type ReportRepo interface {
	GetTotalRevenue(ctx context.Context) (float64, error)
}

type ReportService struct {
	repo ReportRepo
}

func NewReportService(repo ReportRepo) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetRevenue(ctx context.Context) (float64, error) {
	return s.repo.GetTotalRevenue(ctx)
}
