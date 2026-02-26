package service

type AnalyticsRepo interface{ TodaySales() (int64, error) }

type AnalyticsService struct{ repo AnalyticsRepo }

func NewAnalyticsService(repo AnalyticsRepo) *AnalyticsService { return &AnalyticsService{repo: repo} }

func (s *AnalyticsService) TodaySales() (int64, error) { return s.repo.TodaySales() }

func (s *AnalyticsService) WeeklyTrend() ([]int64, error) { return []int64{1, 2, 3}, nil }
