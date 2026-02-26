package service

import "testing"

type fakeAnalyticsRepo struct{}

func (f fakeAnalyticsRepo) TodaySales() (int64, error) { return 1000, nil }

func TestAnalyticsTodaySales(t *testing.T) {
	svc := NewAnalyticsService(fakeAnalyticsRepo{})
	v, _ := svc.TodaySales()
	if v == 0 {
		t.Fatal("expected sales")
	}
}
