package core

import (
	"context"
	"time"
)

type PricePoint struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
}

type TrendService struct{}

func NewTrendService() *TrendService {
	return &TrendService{}
}

func (s *TrendService) GetPriceTrends(ctx context.Context, voyageID string) ([]PricePoint, error) {
	// Mock data
	return []PricePoint{
		{Date: time.Now().AddDate(0, 0, -30).Format("2006-01-02"), Price: 5000},
		{Date: time.Now().AddDate(0, 0, -15).Format("2006-01-02"), Price: 5200},
		{Date: time.Now().Format("2006-01-02"), Price: 5500},
	}, nil
}
